
exports.version = "dejavuDB 0.2.1";
exports.version_info = "";


exports.Get= function(key){
	call_go_fn("Get", key)
}

exports.GetInfo = function(key){

}

exports.Set = function(key, value, type){
	call_go_fn("Set", key, value, type)
}

exports.Update = function(key, value){
	call_go_fn("Update", key, value)
}

exports.contract = class{
	constructor(content){
		if (!(content.constructor == Object)) {
			throw 'expected type "Object"';	
		}

		res = "{";

		for (var [key, value] of Object.entries(content)){
			if (typeof value == 'string'){
				value = '"'+value+'"';
			}
			else{
				value = value.toString();
			}
			res += '"'+value+'":'+value+',';
		}
		
		this.string = res +'}';
	}

}

exports.deployContract = function(key, contract){
	call_go_fn("deployContract",key, contract.string);
}



exports.settings = {

	is_ML_enabled: function(){
		return call_go_fn("dejavu_api_is_ML_enabled");
	},

	enable_ML: function(){
		return call_go_fn("dejavu_api_enable_ML");
	},

	disable_ML: function(){
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