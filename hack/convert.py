import json
import os.path
import re
import ipykernel
import requests

from requests.compat import urljoin
from notebook.notebookapp import list_running_servers
from nbconvert.preprocessors import Preprocessor
from traitlets.config import Config

import nbconvert
import ast
import astunparse
import tarfile

annotation = 'kubeflow/train'

class AnnotationPreprocessor(Preprocessor):
    def check_cell_conditions(self, cell):
        return annotation in cell.source

    def preprocess(self, nb, resources):
        nb.cells = [cell for cell in nb.cells if self.check_cell_conditions(cell)]
        return nb, resources

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


def convert_to_python(notebook_file):
    c = Config()
    c.PythonExporter.preprocessors = [AnnotationPreprocessor]
    exporter = nbconvert.PythonExporter(config=c)
    script, resources = exporter.from_filename(notebook_file)
    return script

def gen_tarball(contents):
    tar_name = "/tmp/output.tar"
    tmpfile = "/tmp/notebook.py"
    with open(tmpfile, "w+") as f:
        f.write(contents)
    with tarfile.open(tar_name, "w:gz") as tar:
        tar.add(tmpfile)
    return tar_name


notebook_name = get_notebook_name()
source = convert_to_python(notebook_name)
print(source)
output_tar = gen_tarball(source)
print(output_tar)

