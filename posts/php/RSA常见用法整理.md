```
{
    "url": "rsa-usage",
    "time": "2014/08/08 11:34",
    "tag": "PHP,RSA"
}
```

RSA是一种非对称加密算法，不但可以支持加解密还支持数字签名；即可以保障数据的安全性也可以保证数据不被篡改。下面整理下常见的用法。

# 私钥签名 - 公钥验签

与其他系统交互之间常常需要确保提交的数据不被篡改，在游戏接入中最常见的做法就是采用MD5加密。双方将通信数据按照一定规则处理并和双方定义好的KEY组合，然后对这串数据进行MD5处理，接受方收到后做同样处理。若一致则表示该请求是正常的，没有经过篡改的；若不一致则认为是异常请求。
该方法简单快捷，但也存在个严重的隐患，就是KEY是双方都知道的，需要共同维护。如果APPKEY不小心被泄漏，知道规则的前提下则可以模拟正常情况。这是存在一定风险的。

RSA签名方式可以很好的解决该问题。RSA签名会生成对应的私钥和公钥，使用私钥对数据加签，然后接收方使用公钥验签。签名的动作只能由私钥进行，验签则由公钥验签。同样的，接收方返回的数据也可以用私钥加签，反过来自己使用对方提供的公钥验签。各自维护各自的私钥，公钥提供给对方验签。可以确保验签通过的数据一定来自对方的私钥签名，从而确保安全性。

PHP中验签和加签操作处理：

```
public static function verify($data, $sign_str, $public_key)
{
    if(empty($sign_str)) {
        return false;
    }
    $res = openssl_get_publickey($public_key);
    if(! is_resource($res)) {
        return false;
    }
    $status = openssl_verify($data, base64_decode($sign_str), $public_key);
    openssl_free_key($res);
    return $status;
}
public static function sign($data, $private_key)
{
    $res = openssl_get_privatekey($private_key);
    if(! is_resource($res)) {
        return false;
    }
    $sign_str = '';
    openssl_sign($data, $sign_str, $res);
    openssl_free_key($res);
    return base64_encode($sign_str);
}
```

# 公钥加密 - 私钥解密

RSA除了签名也可以加解密，和加签、验签一样。常用的公钥加密-私钥解密确保数据只有拥有私钥才可以解密；而私钥加密-公钥解密则可以确保该数据是来自私钥的签名。

与加签验签有一点不同的是，加解密有长度限制，1024位密钥的密文长度不可超过117位，2048位密钥的密文长度不能超过245。

```
public static function pubEncrypt($data, $public_key)
{
    $res = openssl_get_publickey($public_key);
    if(! is_resource($res)) {
        return false;
    }
    if(openssl_public_encrypt($data, $crypt_text, $res)) {
        openssl_free_key($res);
        return base64_encode($crypt_text);
    }
    return false;
}
public static function pubDecrypt($sign, $public_key)
{
    $res = openssl_get_publickey($public_key);
    if(! is_resource($res)) {
        return false;
    }
    if(openssl_public_decrypt(base64_decode($sign), $data, $res)) {
        openssl_free_key($res);
        return $data;
    }
    return false;
}
public static function priEncrypt($data, $private_key)
{
    $res = openssl_get_privatekey($private_key);
    if(! is_resource($res)) {
        return false;
    }
    if(openssl_private_encrypt($data, $crypt_text, $res)) {
        openssl_free_key($res);
        return base64_encode($crypt_text);
    }
    return false;
}
public static function priDecrypt($sign, $private_key)
{
    $res = openssl_get_privatekey($private_key);
    if(! is_resource($res)) {
        return false;
    }
    if(openssl_private_decrypt(base64_decode($sign), $data, $res)) {
        openssl_free_key($res);
        return $data;
    }
    return false;
}
```

# 生成密钥

手册上openssl_pkey_new函数的使用必须安装有效的 openssl.cnf 以保证此函数正确运行。

```
public static function create($bits = 1024)
{
    $res = openssl_pkey_new(array(
        'private_key_bits' => $bits,
        'private_key_type' => OPENSSL_KEYTYPE_RSA 
    ));
    openssl_pkey_export($res, $private_key);
    $public_key = openssl_pkey_get_details($res);
    $public_key = $public_key['key'];
    return array(
        'private_key' => $private_key,
        'public_key' => $public_key
    );
}
```