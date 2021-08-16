package main

import (
	"errors"
	"runtime"
	"sync"
)

var tf_model_register = map[string]*tfModel{}

type tfModel struct {
	name         string
	layer_count  int //keras
	param_count  int //keras
	input_shape  []int
	output_shape []int
	constructer  string //javascript

	path string
	lock sync.Mutex
}

func init_tensorflow() error {
	switch runtime.GOOS {
	case "linux":
	case "darwin":
	case "windows":

	}
	return nil
}

func tf_create_image_classifyer() string {
	return `
	function make_model(input_shape, num_classes){
	inputs = keras.Input(shape=input_shape)

    // Image augmentation block
    x = tf.sequential(
		[
			tf.layers.experimental.preprocessing.RandomFlip("horizontal"),
			tf.layers.experimental.preprocessing.RandomRotation(0.1),
		]
	)(inputs)

    // Entry block
    x = layers.experimental.preprocessing.Rescaling(1.0 / 255)(x)
    x = layers.Conv2D(32, 3, strides=2, padding="same")(x)
    x = layers.BatchNormalization()(x)
    x = layers.Activation("relu")(x)

    x = layers.Conv2D(64, 3, padding="same")(x)
    x = layers.BatchNormalization()(x)
    x = layers.Activation("relu")(x)

    previous_block_activation = x  # Set aside residual

    for size in [128, 256, 512, 728]:
        x = layers.Activation("relu")(x)
        x = layers.SeparableConv2D(size, 3, padding="same")(x)
        x = layers.BatchNormalization()(x)

        x = layers.Activation("relu")(x)
        x = layers.SeparableConv2D(size, 3, padding="same")(x)
        x = layers.BatchNormalization()(x)

        x = layers.MaxPooling2D(3, strides=2, padding="same")(x)

        # Project residual
        residual = layers.Conv2D(size, 1, strides=2, padding="same")(
            previous_block_activation
        )
        x = layers.add([x, residual])  # Add back residual
        previous_block_activation = x  # Set aside next residual

    x = layers.SeparableConv2D(1024, 3, padding="same")(x)
    x = layers.BatchNormalization()(x)
    x = layers.Activation("relu")(x)

    x = layers.GlobalAveragePooling2D()(x)
    if num_classes == 2:
        activation = "sigmoid"
        units = 1
    else:
        activation = "softmax"
        units = num_classes

    x = layers.Dropout(0.5)(x)
    outputs = layers.Dense(units, activation=activation)(x)
    return keras.Model(inputs, outputs)
	}
	const model = make_model();
	model.compile(
		optimizer=keras.optimizers.Adam(1e-3),
		loss="binary_crossentropy",
		metrics=["accuracy"],
	)
	const saveResult = await model.save('http://model-server:5000/upload');
	`
}

func tf_model_predict(model_name string, data interface{}) string {
	return ``
}

func tf_get_model_by_name(name string) (*tfModel, error) {
	if v, ok := tf_model_register[name]; ok {
		return v, nil
	}
	return nil, errors.New("model " + name + " does not exist")
}
