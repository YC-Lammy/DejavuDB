returning_print_buffer = "";

function require(name){
	script = call_go_fn("require", name);
	var exports = {};
	eval(script)
	return exports
}

DB = function(){

	exports = {}

	exports.version = "dejavuDB 0.2.1";
	exports.version_info = "";


	exports.Get= function(key){
		return call_go_fn("Get", key)
	}

	exports.GetInfo = function(key){

	}

	exports.Set = function(key, value, type){
		return call_go_fn("Set", key, value, type)
	}

	exports.Update = function(key, value){
		return call_go_fn("Update", key, value)
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
				res += '"'+String(key)+'":'+value+',';
			}
			
			this.string = res +'}';
			this.IsContract = true
		}

	}

	exports.deployContract = function(key, contract){
		return call_go_fn("deployContract",key, contract.string);
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

	exports.types = {
		string: '\x00',
		String: '\x00'
	}

	exports.Types = function(){
		return{
			string: '\x00',
			String: '\x00'}
	}

return exports

}()

function println(...any){
	for (let i=0; i<any.length; i++) returning_print_buffer += String(any[i])+"\n"
}

function print(...any){
	for (let i=0; i<any.length; i++) returning_print_buffer += String(any[i]);
}