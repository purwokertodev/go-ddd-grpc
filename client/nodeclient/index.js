'use strict';

const grpc = require('grpc');
const employeeProto = grpc.load('employee.proto');

const client = new employeeProto.api.EmployeeService('127.0.0.1:8080', grpc.credentials.createInsecure());

function createEmployee(name, age, address, salary){
  let em = {name: name, age: age, address: address, salary: salary}
  client.createEmployee(em, (err, res) => {
    if (err) {
      console.log(err);
    } else {
      console.log(res);
    }
  });
}

function getEmployee(id){
  let call = client.getEmployee({key: id});
  call.on('data', (em) => {
    console.log(em)
  });
}

getEmployee('9e10c3fc-f3dd-4dca-b7e9-8f1f1b038dcc');

//createEmployee("Andi", 17, "Banjarnegara", 8000.0);
