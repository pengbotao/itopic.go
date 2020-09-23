```
{
    "url": "uploadify-tips",
    "time": "2013/04/21 17:34",
    "tag": "jquery",
    "toc" : "no"
}
```

- 版本：Uploadify Version 3.2
- 官网：http://www.uploadify.com

Uploadify是一款基于Jquery的上传插件，用起来很方便。但上传过程中的提示语言为英文，这里整理下如何修改英文为中文提示。

方法1：直接修改uploadify.js中的提示信息，将英文提示改成对应的中文。不过从软件设计的角度来说，直接修改原类库不是最好的解决方案，会影响到软件的升级。

方法2：重写Uploadify事件`'overrideEvents' : [ 'onDialogClose', 'onUploadError', 'onSelectError' ]`当重写onDialogClose事件后，Uploadify的错误提示信息就都不会提示了。提示信息可直接自定义弹出。

重写事件errorCode的定义在js库中都可以找到。也可以直接用this.queueData.errorMsg来改变提示信息
```
var uploadify_onSelectError = function(file, errorCode, errorMsg) {
        var msgText = "上传失败\n";
        switch (errorCode) {
            case SWFUpload.QUEUE_ERROR.QUEUE_LIMIT_EXCEEDED:
                //this.queueData.errorMsg = "每次最多上传 " + this.settings.queueSizeLimit + "个文件";
                msgText += "每次最多上传 " + this.settings.queueSizeLimit + "个文件";
                break;
            case SWFUpload.QUEUE_ERROR.FILE_EXCEEDS_SIZE_LIMIT:
                msgText += "文件大小超过限制( " + this.settings.fileSizeLimit + " )";
                break;
            case SWFUpload.QUEUE_ERROR.ZERO_BYTE_FILE:
                msgText += "文件大小为0";
                break;
            case SWFUpload.QUEUE_ERROR.INVALID_FILETYPE:
                msgText += "文件格式不正确，仅限 " + this.settings.fileTypeExts;
                break;
            default:
                msgText += "错误代码：" + errorCode + "\n" + errorMsg;
        }
        alert(msgText);
    };
 
var uploadify_onUploadError = function(file, errorCode, errorMsg, errorString) {
        // 手工取消不弹出提示
        if (errorCode == SWFUpload.UPLOAD_ERROR.FILE_CANCELLED
                || errorCode == SWFUpload.UPLOAD_ERROR.UPLOAD_STOPPED) {
            return;
        }
        var msgText = "上传失败\n";
        switch (errorCode) {
            case SWFUpload.UPLOAD_ERROR.HTTP_ERROR:
                msgText += "HTTP 错误\n" + errorMsg;
                break;
            case SWFUpload.UPLOAD_ERROR.MISSING_UPLOAD_URL:
                msgText += "上传文件丢失，请重新上传";
                break;
            case SWFUpload.UPLOAD_ERROR.IO_ERROR:
                msgText += "IO错误";
                break;
            case SWFUpload.UPLOAD_ERROR.SECURITY_ERROR:
                msgText += "安全性错误\n" + errorMsg;
                break;
            case SWFUpload.UPLOAD_ERROR.UPLOAD_LIMIT_EXCEEDED:
                msgText += "每次最多上传 " + this.settings.uploadLimit + "个";
                break;
            case SWFUpload.UPLOAD_ERROR.UPLOAD_FAILED:
                msgText += errorMsg;
                break;
            case SWFUpload.UPLOAD_ERROR.SPECIFIED_FILE_ID_NOT_FOUND:
                msgText += "找不到指定文件，请重新操作";
                break;
            case SWFUpload.UPLOAD_ERROR.FILE_VALIDATION_FAILED:
                msgText += "参数错误";
                break;
            default:
                msgText += "文件:" + file.name + "\n错误码:" + errorCode + "\n"
                        + errorMsg + "\n" + errorString;
        }
        alert(msgText);
    }
    return parameters;
}
 
var uploadify_onSelect = function(){
 
};
 
var uploadify_onUploadSuccess = function(file, data, response) {
    alert(file.name + "\n\n" + response + "\n\n" + data);
};
var uploadify_config = {
    'uploader' : 'upload.php',
    'swf' : '/js/uploadify/uploadify.swf',
    'buttonImage' : '/images/uploadify-button.png',
    'cancelImg' : '/images/uploadify-cancel.png',
    'wmode' : 'transparent',
    'removeTimeout' : 0,
    'width' : 80,
    'height' : 30,
    'multi' : false,
    'auto' : true,
    'buttonText' : '上传',
    'hideButton' : 'true',
    'fileTypeExts' : '*.png;*.jpg;*.jpeg',
    'fileSizeLimit' : '1MB',
    'fileTypeDesc' : 'Image Files',
    'formData' : {"action": "upload", "sid" : ""},
    'overrideEvents' : [ 'onDialogClose', 'onUploadSuccess', 'onUploadError', 'onSelectError' ],
    'onSelect' : uploadify_onSelect,
    'onSelectError' : uploadify_onSelectError,
    'onUploadError' : uploadify_onUploadError,
    'onUploadSuccess' : uploadify_onUploadSuccess
};
 
$("#id").uploadify(uploadify_config);
```
说明：由于FLASH的BUG导致在FF中上传时获取不到SESSION，可以使用formData来传值，如：
```
formData : { '<?php echo session_name();?>' : '<?php echo session_id();?>' }
```