import json
import datetime

from decl_enum import EnumSymbol


class MyJSONEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, datetime.datetime):
            return obj.isoformat()
        elif isinstance(obj, datetime.date):
            return obj.isoformat()
        elif isinstance(obj, datetime.timedelta):
            return (datetime.datetime.min + obj).time().isoformat()
        elif isinstance(obj, EnumSymbol):
            return obj.name
        else:
            return super(MyJSONEncoder, self).default(obj)
