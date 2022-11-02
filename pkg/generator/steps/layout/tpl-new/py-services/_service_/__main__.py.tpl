import {{.ServiceName}}.generated
from {{.ServiceName}}.generated.app.mify_app import MifyServiceApp

from libraries.generated.logs.logger import MifyLoggerWrapper

def main():
    app = MifyServiceApp()
    app.run()

if __name__ == '__main__':
    main()
