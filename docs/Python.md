# Python Lambdas

To work with Python lambdas, you first need to setup your Python virtual environment. To learn more about this please read more about virtual environments and why they are useful [here](https://python.land/virtual-environments/virtualenv). Jeeves requires that you specify a virtual environment because it needs some way to determine where your dependencies are located, and then communicate that to Docker. Because Jeeves is a (near) zero-config tool, checking your virtual environment becomes the only realistic option to pull this off. 

Create a virtual environment

```sh
python3 -m venv [VENV_NAME]
```

Activate the virtual environment

```sh
source [VENV_NAME]/bin/activate
```

Jeeves will also support virtual environments created and managed by `tox`. Just make sure that you activate your environment before running `jeeves faas start`.