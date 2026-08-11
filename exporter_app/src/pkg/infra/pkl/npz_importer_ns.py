import os
import numpy as np
from datetime import datetime


class NpzImporter:

    def __init__(self, directory):
        self.directory = directory

    def __iter__(self):
        return self.import_data()

    def import_data(self):

        files = sorted(
            f for f in os.listdir(self.directory)
            if f.endswith(".npz")
        )


        for file in files:

            arr = np.load(os.path.join(self.directory, file))["arr_0"]

            timestamp = datetime.strptime(
                file.split("_")[1].split(".")[0],
                "%Y%m%d-%H%M%S"
            )

            yield dict({
                "timestamp":timestamp,
                "spectrum":arr
            })