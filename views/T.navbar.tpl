{{define "navbar"}}
<a class="navbar-brand" href="/">北凉刮风了</a>
<div>
        <ul class="nav navbar-nav">
        {{$isLogin := .IsLogin}}
            <li {{if .IsHome}}class="active"{{end}} ><a href="/">首页</a></li>
            {{if $isLogin}}<li {{if .IsCategory}}class="active"{{end}} ><a href="/category">分类</a></li>{{end}}
            <li {{if .IsTopic}}class="active"{{end}} ><a href="/topic">文章</a></li>
            <li  ><a href="/topic">关于</a></li>
        </ul>
</div>
<div class="pull-right">
            <ul class="nav navbar-nav">
                {{if .IsLogin}}
                <li><a href="/login?exit=true">退出</a></li>
                {{else}}
                <li><a href="/login">管理员登录</a></li>
                {{end}}
            </ul>
        </div>
{{end}}
