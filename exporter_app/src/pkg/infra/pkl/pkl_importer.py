import os
import pickle
import re
from datetime import datetime

class PklImporter:
    def __init__(self, pkl_data_path):
        self.path = pkl_data_path

    def __iter__(self):
        return self.import_data()

    # this import only files which is exist when the function is called
    def import_data(self):
        files = os.listdir(self.path)


        def extract_start(f):
            m = re.search(r'temp_(\d{8})-', f)
            if not m:
                raise ValueError(f"invalid filename: {f}")
            return datetime.strptime(m.group(1), "%Y%m%d")

        files = sorted(files, key=extract_start)

        if not files:
            raise FileNotFoundError("No files found in the directory.")

        for file in files:
            file_path = os.path.join(self.path, file)
            try:
                with open(file_path, "rb") as f:
                    df = pickle.load(f)
                    times = df.index.to_list()

                for (_, row), time in zip(df.iterrows(), times):
                    data = row.to_dict()
                    data["TimeStamp"] = time
                    yield data
            except (EOFError, pickle.UnpicklingError, AttributeError, TypeError, FileNotFoundError, IsADirectoryError):
                try:
                    with open(file_path, "r") as f:
                        for line in f:
                            line = line.strip()
                            if line:
                                yield line
                except (FileNotFoundError, IsADirectoryError, PermissionError):
                    continue
