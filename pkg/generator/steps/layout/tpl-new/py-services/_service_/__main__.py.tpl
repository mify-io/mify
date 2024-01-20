# vim: set ft=python:
from {{.Workspace.MifyGeneratedServicePackage}}.app.mify_app import MifyServiceApp


def main():
    app = MifyServiceApp()
    app.run()

if __name__ == '__main__':
    main()
