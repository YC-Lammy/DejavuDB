class dejavu {
	constructor(){

	}
	static version = "dejavuDB 0.2.1";
	static version_info = "";

	static settings = dejavu_api_settings;

	static ML = dejavu_api_ML;


	static Get(key){

	}

    static Set(key, value, type){

	}

	static Update(key, value){

	}

	static Batch(command,args){

	}
}


class dejavu_api_settings {
	constructor(){}

	static is_ML_enabled(){
		return dejavu_api_is_ML_enabled();
	}

	static enable_ML(){
		dejavu_api_enable_ML();
	}

	static disable_ML(){
		dejavu_api_disable_ML();
	}
}



class dejavu_api_ML {
	constructor(){}

	static __name__ = "tensorflow.js";
	static version = "";

	static load_model(name){
		if TF_MODEL_EXIST(name) {
			const model = await tf.loadLayersModel('localhost:7650/'+name);
		}
	}

	static predict(model, value) {
		
	}
}