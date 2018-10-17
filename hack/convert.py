import json
import os.path
import re
import ipykernel
import requests

from requests.compat import urljoin
from notebook.notebookapp import list_running_servers


import nbconvert
import ast
import astunparse

# https://github.com/jupyter/notebook/issues/1000#issuecomment-359875246
def get_notebook_name():
    """
    Return the full path of the jupyter notebook.
    """
    kernel_id = re.search('kernel-(.*).json',
                          ipykernel.connect.get_connection_file()).group(1)
    servers = list_running_servers()
    for ss in servers:
        response = requests.get(urljoin(ss['url'], 'api/sessions'),
                                params={'token': ss.get('token', '')})
        for nn in json.loads(response.text):
            if nn['kernel']['id'] == kernel_id:
                relative_path = nn['notebook']['path']
                return os.path.join(ss['notebook_dir'], relative_path)


def convert_to_python(notebook_file, output_file):
    exporter = nbconvert.PythonExporter()
    script, resources = exporter.from_filename(notebook_file)
    return script

def filter_by_decorator(source, decorator_name):
    tree = ast.parse(source)
    ret = []
    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            decorators = [d.id for d in node.decorator_list]
            if decorator_name in decorators:
                ret.append(astunparse.unparse(node))
    return "\n".join(ret)

notebook_name = get_notebook_name()
source = convert_to_python(notebook_name, "output")
output = filter_by_decorator(source, "kubeflow_model")
print(output)

