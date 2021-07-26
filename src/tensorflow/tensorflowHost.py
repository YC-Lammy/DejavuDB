import tensorflow as tf
import joblib, io, os, socket
from tempfile import NamedTemporaryFile
from pathlib import Path

"""
if tf.config.list_physical_devices("TPU"):
    strategy = tf.distribute.TPUStrategy
elif tf.config.list_physical_devices("GPU"):
    strategy = tf.distribute.MirroredStrategy
else:
    strategy = tf.distribute.get_strategy()
"""
os.chdir(os.path.join(Path.home(),"dejavuDB","ML"))

def load_and_save_model(name:str,modelbytes):
    from tensorflow.keras.models import load_model
    with open(name,'x') as f:
        f.write(modelbytes)
        model = load_model(name)
        return model

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind(("localhost", "3247"))
    s.listen()
    conn, addr = s.accept()