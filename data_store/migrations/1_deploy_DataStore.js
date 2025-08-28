var myContract = artifacts.require("DataStore");   // 把合约名称传给require

module.exports = function(developer){
    // 此处可以deploy多个智能合约
    developer.deploy(myContract);  //部署智能合约时会调用构造函数
}