# gosogouproxy
搜狗浏览器加速代理

本程序将 Sogou 浏览器全网加速功能所使用的 HTTP 代理取出来，单独使用。

使用 Go 语言实现。

目前支持教育网（edu）、电信通（dxt）、网通（cnc）、联通（ctc）和 TC9（疑为铁通）的加速代理，支持 HTTP 的 GET、POST、CONNECT 方法。

此程序的功能的可用性完全依赖于搜狗浏览器的代理协议保持不变，且服务器工作正常。在教育网上使用，该服务于 2014 年曾几度中断。

关于其原理在网上讨论颇多，部分第三方分析如：

* （2009 年）[搜狗浏览器教育网加速所用代理协议初探](http://apt-blog.net/exporing-the-protocol-of-sogou-browser)
* （2011 年）[郁闷！研究了一下Sogou的代理服务器验证协议](http://xiaoxia.org/2011/03/10/depressed-research-about-sogou-proxy-server-authentication-protocol/)
* （2013 年）[搜狗浏览器全网加速](http://zhiwei.li/text/2013/11/%E6%90%9C%E7%8B%97%E6%B5%8F%E8%A7%88%E5%99%A8%E5%85%A8%E7%BD%91%E5%8A%A0%E9%80%9F/)。

目前本程序大致模拟搜狗浏览器 4.1.3.8107 的行为工作。

未实现的功能：

* 从搜狗的服务器上取得代理列表。目前代理列表是直接写在程序里面，但抓包看的结果是代理列表是代理服务器地址列表可以另行取得，但似乎比较麻烦，没有实现。
* 使用 PAC 技术智能选择是否使用代理。搜狗服务器上似乎有混淆过的 PAC 数据，格式不明。 
