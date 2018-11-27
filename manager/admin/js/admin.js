var HGW_MANAGE_API = "//"+window.location.host+"/v1/";
var HGW_MANAGE_DOMAIN = "//"+window.location.host+"/";
function getUrlParam(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
    var r = window.location.search.substr(1).match(reg);  //匹配目标参数
    if (r != null) return unescape(r[2]); return null; //返回参数值
}

function unique(arr){
    var res=[];
    for(var i=0,len=arr.length;i<len;i++){
        var obj = arr[i];
        for(var j=0,jlen = res.length;j<jlen;j++){
            if(res[j]===obj) break;
        }
        if(jlen===j)res.push(obj);
    }
    return res;
}

Array.prototype.remove = function(val) {
    var index = this.indexOf(val);
    if (index > -1) {
        this.splice(index, 1);
    }
};

function isValidIP(ip) {
    var reg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
    return reg.test(ip);
}


function convertByte(limit){
    var size = "";
    if(limit < 0.1 * 1024){                            //小于0.1KB，则转化成B
        size = limit.toFixed(2) + "B"
    }else if(limit < 0.1 * 1024 * 1024){            //小于0.1MB，则转化成KB
        size = (limit/1024).toFixed(2) + "KB"
    }else if(limit < 0.1 * 1024 * 1024 * 1024){        //小于0.1GB，则转化成MB
        size = (limit/(1024 * 1024)).toFixed(2) + "MB"
    }else{                                            //其他转化成GB
        size = (limit/(1024 * 1024 * 1024)).toFixed(2) + "GB"
    }

    var sizeStr = size + "";                        //转成字符串
    var index = sizeStr.indexOf(".");                    //获取小数点处的索引
    var dou = sizeStr.substr(index + 1 ,2)            //获取小数点后两位的值
    if(dou == "00"){                                //判断后两位是否为00，如果是则删除00
        return sizeStr.substring(0, index) + sizeStr.substr(index + 3, 2)
    }
    return size;
}

function request(type, url, dataType, async, data, successCb) {
    $.ajax({
        type: type,
        url: url,
        dataType: dataType,
        async: async,
        data: data,
        success: function(rsp) {
            if(rsp.status == 0) {
                if(rsp.error_code == -2000) {
                    parent.location.href = "init.html";
                } else if(rsp.error_code == -2001) {
                    parent.location.href = "login.html";
                }
            }
            successCb(rsp)
        },
        error: function (data) {

        }
    })
}

function logout() {
    request("GET", HGW_MANAGE_DOMAIN + "logout", "json", false, {},
        function(rsp) {
            if(rsp.status == 1) {
                layer.alert("退出成功", {icon: 6}, function () {
                    parent.location.href = "login.html";
                });
            }
        });
}

