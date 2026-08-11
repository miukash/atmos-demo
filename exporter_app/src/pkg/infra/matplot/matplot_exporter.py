import matplotlib.pyplot as plt

class MatplotlibExporter:

    def export(self, spectrum):

        plt.plot(
            spectrum.frequency,
            spectrum.power,
            label=f"Pixel {spectrum.pixel} time {spectrum.timestamp.strftime('%Y-%m-%d %H:%M:%S')}"
        )