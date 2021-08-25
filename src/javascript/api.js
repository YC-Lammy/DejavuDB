
exports.version = "dejavuDB 0.2.1";
exports.version_info = "";


exports.Get= function(key){

}

exports.GetInfo = function(key){

}

exports.Set = function(key, value, type){

}

exports.Update = function(key, value){

}

exports.Batch= function(command,args){

}



exports.settings = class {

	static is_ML_enabled(){
		return call_go_fn("dejavu_api_is_ML_enabled");
	}

	static enable_ML(){
		return call_go_fn("dejavu_api_enable_ML");
	}

	static disable_ML(){
		return call_go_fn("dejavu_api_disable_ML");
	}
}

exports.ML = class {
	static __name__ = "tensorflow.js";
	static version = "";

	static load_model(name){
	}

    static suspend_model(name){

    }

	static predict(model, value) {
		
	}
}