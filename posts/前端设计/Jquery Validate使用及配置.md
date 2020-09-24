```
{
    "url": "jquery-validate",
    "time": "2013/06/01 20:52",
    "tag": "jquery",
    "toc": "no"
}
```

- 当前版本：v1.11.1
- 官网地址：http://jqueryvalidation.org/

# 导入JS文件
下载压缩包后validate文件位于dist目录，目录中包含jquery.validate.js 与 additional-methods.js以及各自的min文件。additional-methods.js为附加的验证方法，可根据情况导入。
```
<script src="./js/jquery.js" type="text/javascript"></script>
<script src="./js/jquery.validate.js" type="text/javascript"></script>
<script src="./js/additional-methods.js" type="text/javascript"></script>
```

# 使用方法
以注册页为例，需要验证用户名、密码、重复密码、验证码。其中用户名需检测唯一性，验证码需检测是否正确。form-register为表单ID，验证代码如下：
```
$().ready(function(){
    $("#form-register").validate({
        debug: true,
        onfocusout: function (element) {
            $(element).valid();
        },
        errorElement: 'label',
        errorClass: 'text-error',
        submitHandler: function(form) {
            form.submit(); 
        },
        errorPlacement: function(error, element) {
            if(element.attr("name") == 'data[captcha]') {
                error.insertAfter("#captcha_msg");
            } else {
                error.insertAfter(element);
            }
        },
        rules:{
            'data[username]': { required: true, minlength: 6, maxlength:20, lettersonly:true, remote:{
                url: "{{ site_url('ajax/user_check') }}",
                type: "post"
            }},
            'data[password]': { required: true, minlength: 6, maxlength:50},
            'data[repassword]': { required: true, equalTo: "#password"},
            'data[captcha]': { required: true, minlength: 4,remote:{
                url: "{{ site_url('captcha/check') }}",
                type: "post"
            }}
        },
        messages:{}
    });
});
```

# 提示信息
点击提交按钮后验证不通过的会自动在input后增加提示信息
```
<label for="username" class="text-error" style="">This field is required.</label>
```
默认提示信息为英文，可将下面提示信息保存到messages_zh.js并引入。
```
/*
 * Translated default messages for the jQuery validation plugin.
 * Locale: ZH (Chinese, 中文 (Zhōngwén), 汉语, 漢語)
 */
(function ($) {
    $.extend($.validator.messages, {
        required: "必填字段",
        remote: "请修正该字段",
        email: "请输入正确格式的电子邮件",
        url: "请输入合法的网址",
        date: "请输入合法的日期",
        dateISO: "请输入合法的日期 (ISO).",
        number: "请输入合法的数字",
        digits: "只能输入整数",
        creditcard: "请输入合法的信用卡号",
        equalTo: "请再次输入相同的值",
        accept: "请输入拥有合法后缀名的字符串",
        maxlength: $.validator.format("请输入一个长度最多是 {0} 的字符串"),
        minlength: $.validator.format("请输入一个长度最少是 {0} 的字符串"),
        rangelength: $.validator.format("请输入一个长度介于 {0} 和 {1} 之间的字符串"),
        range: $.validator.format("请输入一个介于 {0} 和 {1} 之间的值"),
        max: $.validator.format("请输入一个最大为 {0} 的值"),
        min: $.validator.format("请输入一个最小为 {0} 的值")
    });
}(jQuery));
```
默认的提示信息可能还不够友好，可以进提示信息进行配置，设置messages字段即可
```
messages:{
    'data[username]': {
        required: '请输入用户名',
        minlength: '请输入6-20个英文字符',
        maxlength: '请输入6-20个英文字符',
        lettersonly: '请输入6-20个英文字符',
        remote: '该用户名已被使用'
    },
    'data[password]': {
        required: "请输入密码",
        minlength: jQuery.format("密码不能小于{0}个字 符"),
        maxlength: jQuery.format("密码不能大于{0}个字 符")
    },
    'data[repassword]': {
        required: "请输入确认密码",
        equalTo: "两次密码输入不一致"
    },
    'data[captcha]': {
        required: "请输入验证码",
        minlength: "验证码输入错误",
        remote: "验证码输入错误"
    }
}
```

# 常用设置
- debug：开启调试，当设置true时只验证， 不会提交表单。
- onfocusout：失去焦点验证，上例中是失去焦点就验证，不需要点击submit按钮。
- errorElement: 用来指定错误提示标签，默认为label。
- errorClass: 指定错误提示的css类名,可以自定义错误提示的样式。
- submitHandler：可以接管submit事件。
- errorPlacement：指定错误显示为位置，默认为验证元素后。
- rules：验证规则。
- message：指定提示信息。
- ignore：对某些元素不进行验证

# 自定义验证方法
addMethod(name,method,message)方法：

- 参数name是添加的方法的名字
- 参数method是一个函数,接收三个参数(value,element,param)
    - value是元素的值,
    - element是元素本身
    - param是参数

如additional-methods.js中的lettersonly验证
```
jQuery.validator.addMethod("lettersonly", function(value, element) {
    return this.optional(element) || /^[a-z]+$/i.test(value);
}, "Letters only please");
```