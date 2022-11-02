import sys
from typing import get_type_hints, Union

class MifyConfigError(Exception):
    pass

def _parse_bool(val: Union[str, bool]) -> bool:  # pylint: disable=E1136
    return val if type(val) == bool else val.lower() in ['true', 'yes', '1']

class MifyStaticConfig:
    def __init__(self, env):
        self._environ = env

    def get_config(self, confClass):
        inst = confClass()
        env_map = {}
        if '_env_mapping' in inst.__annotations__:
            env_map = getattr(inst, '_env_mapping', {})
        for field in inst.__annotations__:
            if not field.isupper():
                continue
            env_name = field
            if field in env_map:
                env_name = env_map[field]

            # Raise MifyConfigError if required field not supplied
            default_value = getattr(inst, field, None)
            if default_value is None and self._environ.get(env_name) is None:
                if field != env_name:
                    raise MifyConfigError(f'The {field} (env: {env_name}) field is required')
                else:
                    raise MifyConfigError(f'The {field} field is required')

            # Cast env var value to expected type and raise ValueError on failure
            try:
                var_type = get_type_hints(confClass)[field]
                if var_type == bool:
                    value = _parse_bool(self._environ.get(env_name, default_value))
                else:
                    value = var_type(self._environ.get(env_name, default_value))

                inst.__setattr__(field, value)
            except ValueError:
                raise MifyConfigError('Unable to cast value of "{}" to type "{}" for "{}" field'.format(
                    self._environ[env_name],
                    var_type,
                    field
                )
            )
        return inst
