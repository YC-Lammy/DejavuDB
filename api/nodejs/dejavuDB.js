const Net = require('net');

const client = new Net.Socket();

//const str = new TextDecoder("utf-8").decode;

var connected = false;
var username = "";
var password = "";
var user_login = false;
var buffer_array = [];

client.on("data", function(data) {
    var arrByte = Uint8Array.from(data);
    buffer_array.push(arrByte);
});

client.on('end', function() {
    console.log('Requested an end to the TCP connection');
    connected = false;
});

function send(message) {
    var array = Uint8Array.from(message);
    array.push(0x00);
    client.write(array);
};

function recieve() {
    for (let i = 0; i < buffer_array.length; i++) {
        if (buffer_array[i] == 0x00) {
            result = buffer_array.slice(0, i);
            buffer_array = buffer_array.slice(i + 1, buffer_array.length);
            return result;
        }
    };
};

function connect(host, port) {
    client.connect({ port: port, host: host });
};
exports.connect = connect;

function login(username, password) {
    username = username;
    password = password;
    user_login = true;
};
exports.login = login;

function Set(location, value, type) { // Set location value type
    send("Set " + location + " " + value.toString() + " " + type);
};
exports.Set = Set;

function Update(location, value) { // Set location value
    send("Update " + location + " " + value.toString());
};
exports.Set = Update;

function Delete(target) {
    send("Delete " + target);
};
exports.Clone = Delete;

function Clone(target, destination) {
    send("Clone " + target + " " + destination);
};
exports.Clone = Clone;

function Move(target, destination) {
    send("Move " + target + " " + destination);
};
exports.Clone = Move;