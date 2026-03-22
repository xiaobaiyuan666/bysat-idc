<?php /*a:25:{s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail.tpl";i:1737452904;s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/head.tpl";i:1737452904;s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/menu.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/pageheader.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/breadcrumb.tpl";i:1737452904;s:109:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/hosting.tpl";i:1737452904;s:109:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/general.tpl";i:1737452904;s:102:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/modal.tpl";i:1737452904;s:100:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/pop.tpl";i:1737452904;s:124:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/servicedetail-download.tpl";i:1737452904;s:110:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/cancelrequire.tpl";i:1737452904;s:111:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/dedicated.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/cloud.tpl";i:1737452904;s:102:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/chart.tpl";i:1737452904;s:111:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/zjmfcloud.tpl";i:1737452904;s:110:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/zjmfdcim.tpl";i:1737452904;s:110:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/software.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/tablestyle.tpl";i:1737452904;s:108:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/tablefooter.tpl";i:1737452904;s:105:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/cdn.tpl";i:1737452904;s:105:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/ssl.tpl";i:1737452904;s:105:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/paymodal.tpl";i:1737452904;s:99:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/error/alert.tpl";i:1737452904;s:109:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/upgrade.tpl";i:1737452904;s:123:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/upgrade-configoptions.tpl";i:1737452904;}*/ ?>

<!DOCTYPE html>
<html lang="zh-CN">

<head>
	<meta charset="utf-8" />
	<title><?php echo $Title; ?> | <?php echo $Setting['company_name']; ?></title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta content="<?php echo $Setting['web_seo_desc']; ?>" name="description" />
	<meta content="<?php echo $Setting['web_seo_keywords']; ?>" name="keywords" />
	<meta content="<?php echo $Setting['company_name']; ?>" name="author" />

	<!-- <link rel="shortcut icon" href="/themes/clientarea/default/assets/images/favicon.ico"> -->
<link href="/themes/clientarea/default/assets/css/bootstrap.min.css?v=<?php echo $Ver; ?>" rel="stylesheet" type="text/css" />
<link href="/themes/clientarea/default/assets/css/icons.min.css?v=<?php echo $Ver; ?>" rel="stylesheet" type="text/css" />
<link href="/themes/clientarea/default/assets/css/app.min.css?v=<?php echo $Ver; ?>" rel="stylesheet" type="text/css" />
<?php if(($load_css=load_css('custom.css'))): ?>
    <link href="<?php echo $load_css; ?>?v=<?php echo $Ver; ?>" rel="stylesheet" type="text/css" />
<?php endif; ?>
<!-- 自定义全局样式 -->
<link href="/themes/clientarea/default/assets_custom/css/global.css?v=<?php echo $Ver; ?>" rel="stylesheet" >
<link href="/themes/clientarea/default/assets_custom/css/responsive.css?v=<?php echo $Ver; ?>" rel="stylesheet">
<!-- 字体图标 -->



 <link href="/themes/clientarea/default/assets_custom/fonts/iconfont.css?v=<?php echo $Ver; ?>" rel="stylesheet"> 

<!-- JAVASCRIPT -->
<script src="/themes/clientarea/default/assets/libs/jquery/jquery.min.js?v=<?php echo $Ver; ?>"></script>
<script src="/themes/clientarea/default/assets/libs/bootstrap/js/bootstrap.bundle.min.js?v=<?php echo $Ver; ?>"></script>
<script src="/themes/clientarea/default/assets/libs/metismenu/metisMenu.min.js?v=<?php echo $Ver; ?>"></script>
<script src="/themes/clientarea/default/assets/libs/simplebar/simplebar.min.js?v=<?php echo $Ver; ?>"></script>
<script src="/themes/clientarea/default/assets/libs/node-waves/waves.min.js?v=<?php echo $Ver; ?>"></script>

<!-- <script src="/themes/clientarea/default/assets/libs/error-all/solve-error.js" type="text/javascript"></script> -->
<!-- 自定义js -->
<script src="/themes/clientarea/default/assets_custom/js/throttle.js?v=<?php echo $Ver; ?>"></script>

<link type="text/css" href="/themes/clientarea/default/assets/libs/toastr/build/toastr.min.css?v=<?php echo $Ver; ?>" rel="stylesheet" />
<script src="/themes/clientarea/default/assets/libs/toastr/build/toastr.min.js?v=<?php echo $Ver; ?>"></script>


  <script>
	var setting_web_url = ''
  var language=<?php echo json_encode($_LANG); ?>;
  </script>
	<?php $hooks=hook('client_area_head_output'); if($hooks): foreach($hooks as $item): ?>
			<?php echo $item; ?>
		<?php endforeach; ?>
	<?php endif; ?>
<style>
    .logo-lg img{
      width:150px;
      height:auto;
    }
</style>
</head>
<body data-sidebar="dark">
	<?php if($TplName != 'login' && $TplName != 'register' && $TplName != 'pwreset' && $TplName != 'bind' && $TplName != 'loginaccesstoken'): ?>
	<header id="page-topbar">
		<div class="navbar-header">
			<div class="d-flex">
				<!-- LOGO -->
				<div class="navbar-brand-box">
					<a href="<?php echo $Setting['web_jump_url']; ?>" class="logo logo-dark">
						<?php if($Setting['logo_url_home_mini'] !=''): ?>
						<span class="logo-sm">
							<img src="<?php echo $Setting['logo_url_home_mini']; ?>" alt="" height="32">
						</span>
						<?php endif; ?>
						<span class="logo-lg">
							<img src="<?php echo $Setting['web_logo_home']; ?>" alt="" height="17">
						</span>
					</a>

					<a href="<?php echo $Setting['web_jump_url']; ?>" class="logo logo-light">
						<?php if($Setting['logo_url_home_mini'] !=''): ?>
						<span class="logo-sm" style="overflow: hidden;">
							<img src="<?php echo $Setting['logo_url_home_mini']; ?>" alt="" height="32">
						</span>
						<?php endif; ?>
						<span class="logo-lg">
							<img src="<?php echo $Setting['web_logo_home']; ?>" alt="">
						</span>
					</a>
				</div>

				<button type="button" class="btn btn-sm px-3 font-size-16 header-item waves-effect" id="vertical-menu-btn">
					<i class="fa fa-fw fa-bars"></i>
				</button>


			</div>

			<div class="d-flex">


				<div class="dropdown d-inline-block d-lg-none ml-2 phonehide">
					<button type="button" class="btn header-item noti-icon waves-effect" id="page-header-search-dropdown"
						data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
						<i class="mdi mdi-magnify"></i>
					</button>
					<div class="dropdown-menu dropdown-menu-lg dropdown-menu-right p-0"
						aria-labelledby="page-header-search-dropdown">

						<form class="p-3">
							<div class="form-group m-0">
								<div class="input-group">
									<input type="text" class="form-control" placeholder="Search ..." aria-label="Recipient's username">
									<div class="input-group-append">
										<button class="btn btn-primary" type="submit">
											<i class="mdi mdi-magnify"></i>
										</button>
									</div>
								</div>
							</div>
						</form>
					</div>
				</div>

				<!-- 多语言 -->
				<?php if($Setting['allow_user_language']): ?>
				<div class="dropdown d-inline-block">
					<button type="button" class="btn header-item waves-effect" data-toggle="dropdown" aria-haspopup="true"
						aria-expanded="false">
						<img id="header-lang-img" src="/upload/common/country/<?php echo $LanguageCheck['display_flag']; ?>.png" alt="Header Language" height="16">
					</button>
					<div class="dropdown-menu dropdown-menu-right">
						<!-- wyh 20210329 插件使用 -->
						<?php 
							$parse = parse_url(request()->url());
							$path=$parse['path'];
							$query=$parse['query'];
							$query = preg_replace('/&language=[a-zA-Z0-9_-]+/','',$query);
						 ?>
						<!-- item-->
						<?php if($path=="/addons"): foreach($Language as $key=>$list): ?>
								<a href="?<?php if($query): ?><?php echo $query; ?>&<?php endif; ?>language=<?php echo $key; ?>" class="dropdown-item notify-item language" data-lang="zh-cn">
									<img src="/upload/common/country/<?php echo $list['display_flag']; ?>.png" alt="user-image"
										 class="mr-1" height="12"> <span class="align-middle"><?php echo $list['display_name']; ?></span>
								</a>
							<?php endforeach; else: foreach($Language as $key=>$list): ?>
								<a href="?<?php if($query): ?><?php echo $query; ?>&<?php endif; ?>language=<?php echo $key; ?>" class="dropdown-item notify-item language" data-lang="zh-cn">
									<img src="/upload/common/country/<?php echo $list['display_flag']; ?>.png" alt="user-image"
										 class="mr-1" height="12"> <span class="align-middle"><?php echo $list['display_name']; ?></span>
								</a>
							<?php endforeach; ?>
						<?php endif; ?>

					</div>
				</div>
				<?php endif; ?>
        
				<!-- 购物车 -->
				<div class="dropdown d-none d-lg-inline-block ml-1">
					<button type="button" class="btn header-item noti-icon waves-effect">
						<a href="cart?action=viewcart"><i class="bx bx-cart-alt" style="margin-top: 8px;"></i></a>
							<!-- <?php if(count($CartShopData) != '0'): ?>
							<span class="badge badge-danger badge-pill"><?php echo count($CartShopData); ?></span>
							<?php endif; ?> -->
					</button>
				</div> 

				<!-- 消息 -->
				<div class="dropdown d-none d-lg-inline-block ml-1">
					<a href="message">
						<button type="button" class="btn header-item noti-icon waves-effect">
							<i class="bx bx-bell <?php if($Setting['unread_num']): ?>bx-tada<?php endif; ?>" style="margin-top: 8px;"></i>
							<?php if($Setting['unread_num'] != '0'): ?>
							<span class="badge badge-danger badge-pill"><?php echo $Setting['unread_num']; ?></span>
							<?php endif; ?>
						</button>
					</a>
				</div>

				<!-- 个人中心 -->
				<?php if($Userinfo): ?>
				<div class="dropdown d-inline-block">
					<button type="button" class="btn header-item waves-effect d-inline-flex align-items-center"
						id="page-header-user-dropdown" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
						<div class="user-center_header d-inline-flex align-items-center justify-content-center"
							style="display: inline-block;width: 30px;height: 30px;font-size: 16px;">
							<?php if(preg_match("/^[0-9]*[A-Za-z]+$/is", substr($Userinfo['user']['username'],0,1))): ?> 
							  <?php echo strtoupper(substr($Userinfo['user']['username'],0,1)); elseif(preg_match("/^[\x7f-\xff]*$/", substr($Userinfo['user']['username'],0,3))): ?> 
							  <?php echo substr($Userinfo['user']['username'],0,3); else: ?>
							  <?php echo strtoupper(substr($Userinfo['user']['username'],0,1)); ?> 
							<?php endif; ?>
						</div>
						<span class="d-none d-xl-inline-block ml-1" key="t-henry"><?php echo $Userinfo['user']['username']; ?></span>
						<i class="mdi mdi-chevron-down d-none d-xl-inline-block"></i>
					</button>
					<div class="dropdown-menu dropdown-menu-right">
						<!-- item-->
						<a class="dropdown-item" href="details">
							<i class="bx bxs-user-detail font-size-16 align-middle mr-1"></i>
							<span key="t-profile"><?php echo $Lang['personal_information']; ?></span>
						</a>
						<a class="dropdown-item" href="security">
							<i class="bx bx-cog font-size-16 align-middle mr-1"></i>
							<span key="t-profile"><?php echo $Lang['security_center']; ?></span>
						</a>
						<a class="dropdown-item" href="message">
							<i class="bx bxl-messenger font-size-16 align-middle mr-1"></i>
							<span key="t-profile"><?php echo $Lang['message_center']; ?></span>
						</a>
						<?php if($Setting['certifi_open']==1): ?>
						<a class="dropdown-item" href="verified"> 
							<i class="bx bxs-id-card font-size-16 align-middle mr-1"></i>
							<span key="t-profile"><?php echo $Lang['real_name_authentications']; ?></span>
						</a>
						<?php endif; ?>
						<a class="dropdown-item text-danger" href="logout"><i
								class="bx bx-power-off font-size-16 align-middle mr-1 text-danger"></i> <span
								key="t-logout"><?php echo $Lang['log_out']; ?></span></a>
					</div>
				</div>
				<?php else: ?>
				<div class="pointer d-flex align-items-center">
					<a href="/login" class="text-dark"><?php echo $Lang['please_login']; ?></a>
				</div>
				<?php endif; ?>

			</div>
		</div>
	</header>

	<!-- ========== Left Sidebar Start ========== -->
<?php if($Userinfo): ?>
<div class="vertical-menu">
	<div data-simplebar class="h-100">
		<!--- Sidemenu -->
		<div id="sidebar-menu" class="menu-js">
			<!-- Left Menu Start -->
			<ul class="metismenu list-unstyled" id="side-menu">
			
				<!-- 临时菜单 -->
				<!-- <li>
					<a href="/credit" class="waves-effect">
						<i class="bx bx-home-circle"></i>
						<span>信用额度</span>
					</a>
				</li> -->
				<!-- 临时菜单 -->
				<?php foreach($Nav as $nv): ?>
				<li>
					<a href="<?php if($nv['child']): ?>javascript: ;<?php else: ?><?php echo $nv['url']; ?><?php endif; ?>" class="<?php if($nv['child']): ?>has-arrow<?php endif; ?> waves-effect">
						<?php if($nv['fa_icon']): ?><i class="<?php echo $nv['fa_icon']; ?>"></i><?php endif; if((isset($nv['tag']))): ?>
							<?php echo $nv['tag']; ?>
						<?php endif; ?>
						<span><?php echo $nv['name']; ?></span>
					</a>
					<?php if($nv['child']): ?>
					<ul class="sub-menu mm-collapse" aria-expanded="false">
						<?php foreach($nv['child'] as $subnav): ?>
						<li>
							<a href="<?php if($subnav['child']): ?>javascript: ;<?php else: ?><?php echo $subnav['url']; ?><?php endif; ?>"
								class="<?php if($subnav['child']): ?>has-arrow<?php endif; ?> waves-effect">
								<?php if($subnav['fa_icon']): ?><i class="<?php echo $subnav['fa_icon']; ?>"></i><?php endif; if((isset($subnav['tag']))): ?>
									<?php echo $subnav['tag']; ?>
								<?php endif; ?>
								<span><?php echo $subnav['name']; ?></span>
							</a>
							<?php if($subnav['child']): ?>
							<ul class="sub-menu" aria-expanded="false">
								<?php foreach($subnav['child'] as $submenu): ?>
								<li>
									<a href="<?php if($submenu['child']): ?>javascript: ;<?php else: ?><?php echo $submenu['url']; ?><?php endif; ?>"
										class="<?php if($submenu['child']): ?>has-arrow<?php endif; ?> waves-effect">
										<?php if($submenu['fa_icon']): ?><i class="<?php echo $submenu['fa_icon']; ?>"></i><?php endif; if((isset($submenu['tag']))): ?>
											<?php echo $submenu['tag']; ?>
										<?php endif; ?>
										<span><?php echo $submenu['name']; ?></span>
									</a>
								</li>
								<!-- Nav Level 3 -->
								<?php endforeach; ?>
							</ul>
							<?php endif; ?>
						</li>
						<!-- Nav Level 2 -->
						<?php endforeach; ?>
					</ul>
					<?php endif; ?>
				</li>
				<!-- Nav Level 1 -->
				<?php endforeach; ?>
			</ul>
		</div>
		<!-- Sidebar -->
	</div>
</div>
<?php else: ?>
<div class="vertical-menu menu-js">
	<div data-simplebar class="h-100">
		<!--- Sidemenu -->
		<div id="sidebar-menu" class="menu-js">
			<!-- Left Menu Start -->
			<ul class="metismenu list-unstyled" id="side-menu">
				<li>
					<a href="/clientarea" class="waves-effect">
						<i class="bx bx-home-circle"></i>
						<span>首页</span>
					</a>
				</li>
				<li>
					<a href="/login" class="waves-effect">
						<i class="bx bx-user"></i>
						<span>登录</span>
					</a>
				</li>
				<li>
					<a href="/register" class="waves-effect">
						<i class="bx bx-user"></i> 
						<span>注册</span>
					</a>
				</li>
				<li>
					<a href="/cart" class="waves-effect">
						<i class="bx bx-cart-alt"></i>
						<span>订购产品</span>
					</a>
				</li>
				<li>				
					<a href="/news" class="waves-effect">
						<i class="bx bx-detail"></i>
						<span>新闻中心</span>
					</a>
				</li>
				<li>				
					<a href="/knowledgebase" class="waves-effect">
						<i class="bx bx-detail"></i>
						<span>帮助中心</span>
					</a>
				</li>
				<li>				
					<a href="/downloads" class="waves-effect">
						<i class="bx bx-download"></i>
						<span>资源下载</span>
					</a>
				</li>
			</ul>
		</div>
		<!-- Sidebar -->
	</div>
</div>
<?php endif; ?>


	<div class="main-content">
		<div class="page-content">
			<?php if($TplName != 'clientarea'): ?>
			
<div class="container-fluid">
    <!-- start page title -->
    <div class="row">
        <div class="col-12">
            <div class="page-title-box d-flex align-items-center justify-content-between">
                <?php if($TplName == 'viewbilling'): ?>
                <h4 class="mb-0 font-size-18"><?php echo $Title; ?> - <?php echo $Get['id']; ?></h4>
                <?php else: ?>
                <div style="display:flex;">
                    
                    <a href="javascript:history.go(-1)" class="backBtn" style="display: none;"><i class="bx bx-chevron-left" style="font-size: 32px;margin-top: 1px;color: #555b6d;"></i></a>
                    <h4 class="mb-0 font-size-18"><?php echo $Title; ?></h4>
                </div>
                <?php endif; ?>
                <div class="page-title-right">
	                <?php if(!$ShowBreadcrumb): ?>
                    <ol class="breadcrumb m-0">
    <li class="breadcrumb-item"><a href="javascript: void(0);"><?php echo $Lang['title_clientarea']; ?></a></li>
    <li class="breadcrumb-item active"><?php echo $Title; ?></li>
</ol>
                    <?php endif; ?>
                </div>

            </div>
        </div>
    </div>
    <!-- end page title -->    
</div>
<script>
    $(function(){
        $('.backBtn').hide()
    })
</script>
			<?php endif; ?>
			<div class="container-fluid">
				<?php endif; ?><style>
    .error-tip{
        color: #f46a6a;
        margin: 0;
        padding: 0;
        line-height: 36px;
        margin-left:13rem;
        display: none;
    }
    .ml-8{
        margin-left:8.3rem
    }
    .contract_mc{
        width: 100%;
        height: 100%;
        position: fixed;
        right: 0px;
        top: 0px;
        background: #000;
        z-index: 999999;
        opacity: 0.4;
    }
    .pt-9{
        padding-top:9px;
    }
    .must-reinstall-check:before{
        content: '*';
        color: red;
    }
    .d-flex-cl{
        flex-direction: column;
    }
    .d-flex-cl .reinstallAgreeCheckbox{
        color: #ff0000;
    }
</style>
<script src="/themes/clientarea/default/assets/libs/moment/moment.js?v=<?php echo $Ver; ?>"></script>
<script src="/themes/clientarea/default/assets/libs/clipboard/clipboard.min.js?v=<?php echo $Ver; ?>"></script>
<!-- select -->
<link rel="stylesheet" href="/themes/clientarea/default/assets/libs/bootstrap-select/css/bootstrap-select.min.css?v=<?php echo $Ver; ?>">
<script src="/themes/clientarea/default/assets/libs/bootstrap-select/js/bootstrap-select.min.js?v=<?php echo $Ver; ?>"></script>
<script src="/themes/clientarea/default/assets/js/public.js?v=<?php echo $Ver; ?>"></script>
<!-- <link rel="stylesheet" href="<?php echo app('request')->domain(); ?><?php echo app('request')->rootUrl(); ?>/vendor/dcimcloud/css/layui.css"> 隐藏layui把用到layui样式的地方改成bootstrap的样式 -->
<link rel="stylesheet" href="<?php echo app('request')->domain(); ?><?php echo app('request')->rootUrl(); ?>/vendor/dcimcloud/css/selectFilter.css">
<script src="<?php echo app('request')->domain(); ?><?php echo app('request')->rootUrl(); ?>/vendor/dcimcloud/js/sweetAlert2.min.js"></script>
<script src="<?php echo app('request')->domain(); ?><?php echo app('request')->rootUrl(); ?>/vendor/dcimcloud/js/selectFilter.js"></script>
<script src="/themes/clientarea/default/assets/libs/echarts/echarts.min.js?v=<?php echo $Ver; ?>"></script>
<?php if($Detail['host_data']['type'] == "hostingaccount"): ?>
    
<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<div class="container-fluid">
  <div class="row">
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body">
          <div class="row">
            <div class="col-lg-5 mb-1">
              <div class="p-5 bg-primary rounded text-white d-flex flex-column
                                                justify-content-center align-items-center">
                <h1 class="text-white"><?php echo $Detail['host_data']['productname']; ?></h1>
                <p class="mb-4"><?php echo $Detail['host_data']['domain']; ?></p>
                <p>
                  <!-- 备注 -->
                  <span class="text-white-50"><?php echo $Lang['remarks_infors']; ?>： <?php if($Detail['host_data']['remark']): ?><?php echo $Detail['host_data']['remark']; else: ?>-<?php endif; ?></span>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>
                <span class="badge badge-pill py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>
                  mb-3">
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>
                </span>
              </div>
            </div>
            <div class="col-lg-7 mb-1">
              <div class="d-flex justify-content-between">
                <div class="table-responsive" style="min-height: auto;">
                  <table class="table mb-0 table-bordered">
                    <tbody>
                      <tr>
                        <th scope="row"><?php echo $Lang['price']; ?></th>
                        <td>
                          <?php echo $Detail['host_data']['firstpaymentamount_desc']; ?>
                          <span
                            class="ml-2 badge <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): else: ?> <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span>
                        </td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['subscription_date']; ?></th>
                        <td><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['payment_cycle']; ?></th>
                        <td><?php echo $Detail['host_data']['billingcycle_desc']; ?></td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['due_date']; ?></th>
                        <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                        <td>
                          <?php if($Detail['host_data']['billingcycle_desc'] == '一次性' || $Detail['host_data']['billingcycle_desc'] == '免费'): ?>
                            -
                          <?php else: ?>
                            <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?>
                          <?php endif; ?>
                        </td>
                        <?php else: ?>
                        <td><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></td>
                        <?php endif; ?>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['automatic_balance_renewal']; ?></th>
                        <td>
						<?php if($Detail['host_data']['billingcycle'] != 'onetime' && $Detail['host_data']['status'] == 'Paid' && $Detail['host_data']['billingcycle'] != 'free'): ?>
                          <div class="custom-control custom-switch custom-switch-md mb-3" dir="ltr">
                      <input type="checkbox" <?php if($Detail['host_data']['billingcycle_desc'] == '一次性' || $Detail['host_data']['billingcycle_desc'] == '免费'): ?> disabled <?php endif; ?> class="custom-control-input" id="automaticRenewal"
                              onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew']
                              !=0): ?>checked <?php endif; ?>> <label class="custom-control-label" for="automaticRenewal"></label>
                          </div>
						  <?php else: ?>
						  -
						  <?php endif; ?>
                        </td>
                      </tr>
                    </tbody>
                  </table>

                </div>

              </div>


              <div>
                
                <button type="button" class="btn btn-primary" id="logininfo">
                  <?php echo $Lang['login_information']; ?>

                  <i class="mdi mdi-chevron-down"></i>
                </button>

                <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): ?>
                <button type="button" class="btn btn-primary waves-effect waves-light" id="renew"
                  onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>
                <span>
                  <?php if($Cancel['host_cancel']): ?>
                  <button class="btn btn-danger" id="cancelStopBtn"
                    onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['stop_when_due']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-primary" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>
                  <?php endif; ?>
                </span>
                <?php endif; ?>
                <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.3-->
                <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
                    <div class="btn-group ml-0 mr-2">
                      <span>
                      <?php echo $item; ?>
                      </span>
                    </div>
                  <?php endforeach; ?>
                <?php endif; ?>
                <!-- 结束 -->

                <?php if($Detail['module_button']['control']): ?>

                <div class="btn-group">

                  <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i
                      class="mdi mdi-chevron-down"></i></button>

                  <div class="dropdown-menu">

                    <?php foreach($Detail['module_button']['control'] as $item): ?>

                    <a class="dropdown-item service_module_button" href="javascript:void(0);"
                      onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                      data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                      data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

                    <?php endforeach; ?>

                  </div>

                </div>
                <?php endif; if($Detail['module_button']['console']): ?>

                <div class="btn-group">
                  <?php if(($Detail['module_button']['console']|count) == 1): foreach($Detail['module_button']['console'] as $item): ?>
                  <a class="btn btn-primary service_module_button d-flex align-items-center" href="javascript:void(0);"
                    onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                    data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
                  <?php endforeach; else: ?>
                  <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i
                      class="mdi mdi-chevron-down"></i></button>

                  <div class="dropdown-menu">

                    <?php foreach($Detail['module_button']['console'] as $item): ?>

                    <a class="dropdown-item service_module_button" href="javascript:void(0);"
                      onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                      data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

                    <?php endforeach; ?>

                  </div>
                  <?php endif; ?>
                </div>
                <?php endif; ?>

              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body" style="min-height: 500px;">
          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">
            <?php foreach($Detail['module_client_area'] as $key=>$item): ?>

            <li class="nav-item">

              <a class="nav-link <?php if($key==0): ?>active<?php endif; ?>" data-toggle="tab"
                href="#module_client_area_<?php echo $item['key']; ?>" role="tab">
                <span><?php echo $item['name']; ?></span>

              </a>

            </li>
            <?php endforeach; if($Detail['config_options']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area']): ?>active<?php endif; ?>" data-toggle="tab" href="#profile1"
                role="tab">
                <span><?php echo $Lang['configuration_option']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area'] && !$Detail['config_options']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#downgrade" role="tab">
                <span><?php echo $Lang['upgrade_downgrade']; ?></span>
              </a>
            </li>
            <?php endif; ?>

            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area'] && !$Detail['config_options'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#finance" role="tab">
                <span><?php echo $Lang['financial_information']; ?></span>
              </a>
            </li>
            <?php if($Detail['download_data']): ?>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#download" role="tab">
                <span><?php echo $Lang['file_download']; ?></span>
              </a>
            </li>
            <?php endif; ?>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#settings1" role="tab">
                <span><?php echo $Lang['journal']; ?></span>
              </a>
            </li>

          </ul>

          <!-- Tab panes -->
          <div class="tab-content p-3 text-muted">
			<?php if($Detail['config_options']): ?>
            <div class="tab-pane  <?php if(!$Detail['module_client_area']): ?>active<?php endif; ?>" id="profile1" role="tabpanel">
              <div class="row">
                <?php foreach($Detail['config_options'] as $item): ?>
                <div class="col-md-2 mb-2">
                  <div class="bg-light">
                    <div class="card-body">
                      <p><?php echo $item['name']; ?></p>
                      <span><?php echo $item['sub_name']; ?></span>
                    </div>
                  </div>
                </div>
                <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
                <div class="col-md-2 mb-2">
                  <div class="bg-light">
                    <div class="card-body">
                      <p><?php echo $item['fieldname']; ?></p>
                      <span><?php echo $item['value']; ?></span>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
                <?php endforeach; ?>

              </div>
            </div>
			<?php endif; ?>
            <div class="tab-pane <?php if(!$Detail['module_client_area'] && !$Detail['config_options'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>" id="finance" role="tabpanel">

            </div>
            <div class="tab-pane" id="settings1" role="tabpanel">

            </div>
            <?php if($Detail['download_data']): ?>
            <div class="tab-pane" id="download" role="tabpanel">
              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
            </div>
            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <div class="tab-pane" id="downgrade" role="tabpanel">
              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-md-3">

                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>
                        <div class="col-md-6">
                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>
                        </div>
                        <div class="col-md-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-md-3">
                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                        </div>
                        <div class="col-md-6">
                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>
                        </div>
                        <div class="col-md-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; foreach($Detail['module_client_area'] as $key=>$item): ?>

            <div class="tab-pane <?php if($key==0): ?>active<?php endif; ?>" role="tabpanel"
              id="module_client_area_<?php echo $item['key']; ?>">

              <div class="width:100%">
                <script>
                  $.ajax({
                    url : '/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>&date='+Date.parse(new Date()) 
                    ,type : 'get'
                    ,success : function(res) {
                      $('#module_client_area_<?php echo $item['key']; ?> > div').html(res);
                    }
                  })
                </script>
              </div>

              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                onload="this.height=$($('.main-content .card-body')[1]).height()-72" frameborder="0"
                width="100%"></iframe> -->

            </div>

            <?php endforeach; ?>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="deactivateDia" style="display: none;">
    <form>
      <input type="hidden" value="<?php echo $Token; ?>" />
      <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>" />
      <div class="form-group row mb-4">
        <label class="col-sm-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
        <div class="col-sm-8">
          <select class="form-control" class="second_type" name="type">
            <option value="Immediate"><?php echo $Lang['remarks_infors']; ?>立即</option>
            <option value="Endofbilling" selected><?php echo $Lang['billing_cycle']; ?></option>
          </select>
        </div>
      </div>
      <div class="form-group row mb-0">
        <label class="col-sm-3 col-form-label text-right"><?php echo $Lang['cancelreason']; ?></label>
        <div class="col-sm-8">
          <div class="input-group">
            <select class="form-control" class="second_type" name="reason">
              <?php foreach($Detail['cancelist'] as $item): ?>
              <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
              <?php endforeach; ?>
            </select>
          </div>
        </div>
    </form>
  </div>
</div>

<style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<script src="/themes/clientarea/default/assets/libs/clipboard/clipboard.min.js?v=<?php echo $Ver; ?>"></script>
<script>
  function refresh(type) {
    location.reload();
  }


  // 查看密码
  var showPWd = false
  $('#copyPwdContent').hide()
  function togglePwd() {
    showPWd = !showPWd

    if (showPWd) {
      $('#copyPwdContent').show()
      $('#hidePwdBox').hide()
    }
    if (!showPWd) {
      $('#copyPwdContent').hide()
      $('#hidePwdBox').show()
    }
  }

  // 复制密码
  var clipboard = new ClipboardJS('.btn-copy', {
    text: function (trigger) {
      return $('#copyPwdContent').text()
    }
  });
  clipboard.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

</script>
<script>
  const logObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'log_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: logObj,
    success: function (data) {
      $(data).appendTo('#settings1');
    }
  });
  const financialObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'billing_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: financialObj,
    success: function (data) {
      $(data).appendTo('#finance');
    }
  });

  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      
    `)
  });


  $('#popModal').on('shown.bs.modal',function() {
    if (clipboardpoppwd) {
        clipboardpoppwd.destroy()
      }
     clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })
</script>
<?php elseif($Detail['host_data']['type'] == "server"): ?>
    <style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<style>
  .server_header_box {
    height: auto;
    background-image: linear-gradient(87deg, #4d83ff 0%, #3656ff 100%);
    border-radius: 15px;
    padding: 20px 25px;
    color: #ffffff;
  }

  .left_wrap_btn {
    display: inline-block;
    width: 80px;
    height: 20px;
    background-color: #5f88fe;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    color: #ffffff;
    text-align: center;
    border: none;
  }

  .custom-button {
    background-color: #6f87fc;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    font-size: 12px;
    color: #fff;
    border: none;
  }

  .box_left_wrap {
    border-left: 1px solid rgba(255, 255, 255, 0.25);
    min-height: 74px;
  }

  .aibiao a {
    width: 100%;
    height: 100%;
    display: inline-block;
  }

  @media screen and (max-width: 1367px) {
    .form-control {
      width: 46%;
    }

    .server_header_box {
      height: auto;
    }

    .power_box {
      max-width: 300px;
    }

    .left_wrap_btn {
      width: 60px !important;
    }

    .bottom-box {
      margin-top: 3rem !important;
    }

    .osbox {
      max-width: 150px;
    }
  }

  @media screen and (max-width: 976px) {
    .server_header_box {
      height: auto;
      padding: 20px;
      margin-top: 10px;
    }

    .domain,
    .box_left_wrap {
      margin-bottom: 20px;
      border-left: none;
    }

    .power_box {
      margin-bottom: 20px;
    }
  }

  .tuxian {
    cursor: pointer;
  }

  .tuxian:hover {
    color: rgba(224, 224, 224, 0.877);
  }

  .alarm {
    display: inline-block;
    font-size: 12px;
    cursor: pointer;
    color: #495057;
    font-weight: 300;
  }

  .fr {
    float: right;
  }

  .restall-btn {
    border-radius: 25px;
    margin-left: 20px;
  }

  .login-info-icon {
    color: #5f88fe;
  }

  .dc {
    color: #5f88fe
  }

  .rsb {
    height: 20px;
    padding: 0px 10px;
  }

  .mg-0 {
    margin: 0;
  }

  .plr-0 {
    padding-left: 0px;
    padding-right: 0px;
    margin-bottom: 0;
  }

  #copyIPContent:hover {
    color: #FCA426
  }

  #copyOneIp:hover {
    color: #FCA426
  }

  .text-nowrap {
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .text-right {
    text-align: right;
  }

  .pre-money-box {
    background: url("/themes/clientarea/default/assets/images/money.png") no-repeat;
    background-position-x: right;
    background-position-y: bottom;
  }

  .w-75 {
    width: 75% !important;
  }

  .ll-flex {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .xf-bg {
    display: none;
    position: absolute;
    background: #fff;
    padding: 10px 15px;
    border-radius: 4px;
    top: -40px;
    left: 0px;
    box-shadow: 0px 3px 5px 0px rgba(0, 28, 144, 0.21);
    font-size: 12px;
  }

  .xf-bg-text {
    color: #333;
    word-break: break-all;
  }

  .flex-wrap {
    display: flex !important;
    flex-flow: wrap;
  }

  .configuration-btn-down {
    width: 100%;

    text-align: center;
    line-height: 36px;
    color: #5F88FE;
  }

  .configuration-btn-down::after {
    content: '';
    display: inline-block;
    width: 8px;
    height: 8px;
    margin-left: 8px;
    background-color: transparent;
    transform: rotate(225deg);
    border: 1px solid #5F88FE;
    border-bottom: none;
    border-right: none;
    transform-origin: 2px;
    transition: all .2s;
  }

  .configuration-btn-down.isClick::after {
    transform: rotate(45deg);
  }
</style>
<script src="/themes/clientarea/default/assets/libs/echarts/echarts.min.js?v=<?php echo $Ver; ?>"></script>
<div class="container-fluid">
  <div class="row mb-4">
    <div class="col-12">
      <div class="row align-items-center server_header_box">

        <div class="mr-3 power_box">

          <div class="text-white d-flex">
            <!-- 电源状态 -->
            <div class="mr-3 pointer">
              <?php if($Detail['module_power_status'] == '1'): ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="powerBox">
                <span id="powerStatusIcon" class="bx bx-loader" data-toggle="popover" data-trigger="hover" title=""
                  data-html="true" data-content="<?php echo $Lang['loading']; ?>..."></span>
              </div>
              <?php else: ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="statusBox"></div>
              <?php endif; ?>
            </div>
            <div>
              <section class="d-flex align-items-center mb-2">

                <h4 class="text-white mb-0 font-weight-bold"><?php echo $Detail['host_data']['productname']; ?></h4>

                <span class="badge badge-pill ml-2 py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>"
                  style="position: relative;">
                  <div class="xf-bg">
                    <div class="xf-bg-text"><span style="color: #e31519;">暂停原因：</span><?php echo !empty($Detail['host_data']['suspendreason']) ? $Detail['host_data']['suspendreason'] : $Detail['host_data']['suspendreason_type']; ?></div>
                    <font class="sj"></font>
                  </div>
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>
                </span>

              </section>

              <section>

                <span><?php echo $Detail['host_data']['domain']; ?></span>
                <span class="cancelBtn" id="cancelDcimTask" style="display:none;"
                  onclick="cancelDcimTask('<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['cancel_task']; ?></span>
              </section>
            </div>

          </div>

        </div>

        <div class="pl-4 mr-3 box_left_wrap osbox">

          <span class="text-white-50 fs-12"><?php echo $Lang['operating_system']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $Detail['host_data']['os']; ?></h5>

          <?php if(in_array('reinstall', array_column($Detail['module_button']['control'], 'func')) &&
          $Detail['host_data']['domainstatus']=="Active"): ?>

          <span class="ml-0">

            <button type="button" class="service_module_button left_wrap_btn fs-12 restall-btn" data-func="reinstall"
              data-type="default"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reinstall_system']; ?></button>

          </span>

          <?php endif; ?>

        </div>

        <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '6'||$item['option_type'] == '8'): ?>

        <div class="pl-4 mr-3 box_left_wrap">

          <span class="text-white-50 fs-12"><?php echo $item['name']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $item['sub_name']; ?></h5>

        </div>

        <?php endif; ?>

        <?php endforeach; ?>

        <div class="pl-4 mr-3 box_left_wrap">
          <span class="text-white-50 fs-12"><?php echo $Lang['ip_address']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white">
            <!-- <span data-toggle="popover" data-trigger="hover" title="" data-html="true" data-content="
          <?php foreach($Detail['host_data']['assignedips'] as $list): ?>
          <div><?php echo $list; ?></div>
          <?php endforeach; ?>
        "> -->
            <span>
              <?php if($Detail['host_data']['dedicatedip']): if($Detail['host_data']['assignedips']): ?>
              <span class="tuxian" id="copyIPContent"
                class="pointer"><?php echo $Detail['host_data']['dedicatedip']; ?>(<?php echo count($Detail['host_data']['assignedips']); ?>)</span>
              <!-- <span>(<?php echo count($Detail['host_data']['assignedips']); ?>)</span> -->
              <?php else: ?>
              <span id="copyOneIp" class="pointer copyOneIp"><?php echo $Detail['host_data']['dedicatedip']; ?></span>
              <?php endif; ?>
              <!-- <i class="bx bx-copy pointer text-white ml-1 btn-copy" id="btnCopyIP" data-clipboard-action="copy"
          data-clipboard-target="#copyIPContent"></i> -->
              <?php else: ?>
              -
              <?php endif; ?>

            </span>

          </h5>



        </div>

        <!--
<div class="pl-4 mr-3 box_left_wrap">

  <span class="text-white-50 fs-12"><?php echo $Lang['password']; ?></span>
  <h5 class="mt-2" data-toggle="popover" data-trigger="hover" data-html="true"
    data-content="<?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?><br><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?>">
    <span id="hidePwdBox" class="text-white">***********</span>
    <span id="copyPwdContent" class="text-white"><?php echo $Detail['host_data']['password']; ?></span>
    <i class="fas fa-eye pointer ml-2 text-white" onclick="togglePwd()"></i>
    <i class="bx bx-copy pointer ml-1 btn-copy text-white" id="btnCopyPwd" data-clipboard-action="copy"
      data-clipboard-target="#copyPwdContent"></i>
  </h5>

  

  <?php if(in_array('crack_pass', array_column($Detail['module_button']['control'], 'func')) && $Detail['host_data']['domainstatus']=="Active"): ?>

  <button type="button" class="service_module_button ml-0 left_wrap_btn fs-12" data-func="crack_pass"
    data-type="default"
    onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reset_password']; ?></button>

  <?php endif; ?>

</div>
-->

        <div class="d-flex justify-content-end flex-shrink-1 flex-grow-1">

          <?php if($Detail['module_button']['control'] && $Detail['host_data']['domainstatus']=="Active"): ?>

          <div class="btn-group ml-2 mr-2 mt-2">

            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['control'] as $item): if($item['func'] != 'crack_pass' && $item['func'] != 'reinstall'): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

              <?php endif; ?>

              <?php endforeach; ?>

            </div>

          </div>

          <?php endif; if($Detail['module_button']['console'] && $Detail['host_data']['domainstatus']=="Active"): ?>

          <div class="btn-group ml-2 mr-2 mt-2">
            <?php if(($Detail['module_button']['console']|count) == 1): foreach($Detail['module_button']['console'] as $item): ?>
            <a class="btn btn-primary service_module_button d-flex align-items-center custom-button"
              href="javascript:void(0);"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
              data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
            <?php endforeach; else: ?>
            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['console'] as $item): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

              <?php endforeach; ?>

            </div>
            <?php endif; ?>

          </div>

          <?php endif; ?>

          <!-- <div class="btn-group ml-2 mr-2 mt-2">

    <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>

    <span>

      <?php if($Cancel['host_cancel']): ?>

      <button class="btn btn-danger mb-1 h-100" id="cancelStopBtn"
        onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
        ==
        'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button>

      <?php else: ?>

      <button class="btn btn-primary mb-1 h-100 custom-button" data-toggle="modal"
        data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>

      <?php endif; ?>

    </span>

    <?php endif; ?>
  </div> -->

          <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.2-->
          <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
          <div class="btn-group ml-2 mr-2 mt-2">
            <span>
              <?php echo $item; ?>
            </span>
          </div>
          <?php endforeach; ?>
          <?php endif; ?>
          <!-- 结束 -->
        </div>

      </div>
    </div>
  </div>
  <div class="row">
    <div class="col-md-3">
      <div class="card">
        <div class="card-body">
          <!-- <div class="mb-3 text-center" id="logininfo" style="width: 100px;height: 30px;line-height: 30px;background-color: #ffffff;box-shadow: 0px 4px 20px 2px rgba(6, 75, 179, 0.08);border-radius: 4px;cursor: pointer;">
            <?php echo $Lang['login_information']; ?>
          </div> -->

          <!-- 登录信息start -->
          <div class="row">
            <?php if($Detail['host_data']['domainstatus'] == 'Active'): ?>
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['login_information']; ?>
                  <i class="bx bx-user login-info-icon"></i>
                </p>
                <p class="mb-0"><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></p>
                <p class="mb-0"><?php echo $Lang['password']; ?>：
                  <!-- <span id="hidePwdBox" class="text-black">***********</span> -->
                  <?php if($Detail['host_data']['password'] == ''): ?>
                  <span class="text-black btnCopyPwd pointer dc">-</span>
                  <?php else: ?>
                  <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制"
                    id="copyPwdContent" class="text-black btnCopyPwd pointer dc"><?php echo $Detail['host_data']['password']; ?></span>
                  <?php endif; ?>
                  <!-- <i class="fas fa-eye pointer ml-2 text-black" onclick="togglePwd()"></i> -->
                  <!-- <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPwd" data-clipboard-action="copy"
                    data-clipboard-target="#copyPwdContent"></i> -->

                </p>
                <p class="mb-0"><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] ==
                  '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></p>

              </div>
            </div>
            <?php endif; if(($temp_custom_field_data = array_column($Detail['custom_field_data'], 'value', 'fieldname')) &&
            (isset($temp_custom_field_data['panel_address']) || isset($temp_custom_field_data['面板管理地址']) ||
            isset($temp_custom_field_data['panel_passwd']) || isset($temp_custom_field_data['面板管理密码']))): ?>
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['panel_manage_info']; ?>
                  <i class="bx bx-receipt dc"></i>
                </p>
                <?php if(isset($temp_custom_field_data['panel_address']) || isset($temp_custom_field_data['面板管理地址'])): ?>
                <p class="mb-0"><?php echo $Lang['panel_manage_address']; ?>：<?php echo isset($temp_custom_field_data['panel_address']) ? $temp_custom_field_data['panel_address'] : $temp_custom_field_data['面板管理地址']; ?></p>
                <?php endif; if(isset($temp_custom_field_data['panel_passwd']) || isset($temp_custom_field_data['面板管理密码'])): ?>
                <!-- <p class="mb-0"><?php echo $Lang['panel_manage_password']; ?>：<span id="hidePanelPasswd">***********</span> -->
                <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制" id="panelPasswd"
                  class="btnCopyPanelPasswd dc pointer"><?php echo isset($temp_custom_field_data['panel_passwd']) ? $temp_custom_field_data['panel_passwd'] : ($temp_custom_field_data['面板管理密码'] ?:'--'); ?></span>
                <!-- <i class="fas fa-eye pointer ml-2 text-black" onclick="togglePanelPasswd()"></i>
                  <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPanelPasswd" data-clipboard-action="copy" data-clipboard-target="#panelPasswd"></i> -->
                </p>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; ?>
          </div>
          <!-- 登录信息end -->

          <!-- 流量 -->
          <?php if(($Detail['host_data']['domainstatus'] == 'Active' || ($Detail['host_data']['domainstatus'] == 'Suspended' && $Detail['host_data']['suspendreason_type'] == 'flow')) && $Detail['host_data']['bwlimit'] > 0): ?>
            <!-- <div class="d-flex justify-content-end mb-2">
            <button type="button" class="btn btn-success btn-sm waves-effect waves-light"
              id="orderFlowBtn" onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
          </div> -->
          <div class="mt-4 mb-3">
            <i class="bx bx-circle" style="color:#f0ad4e"></i> <?php echo $Lang['used_flow']; ?>：<span id="usedFlowSpan">-</span>
            <i class="bx bx-circle" style="color:#34c38f"></i> <?php echo $Lang['residual_flow']; ?>：<span id="remainingFlow">-</span>
          </div>
          <div class="mb-4 ll-flex">
            <div class="progress w-75">
              <div class="progress-bar progress-bar-striped bg-success" id="totalProgress" role="progressbar"
                style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">100%
              </div>
            </div>
            <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb" id="orderFlowBtn"
              onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
          </div>
          <?php endif; ?>
          <!-- 流量end -->


          <!-- 订购价格start -->
          <div class="row">
            <!-- <div class="col-12 my-2">
              <div class="d-flex justify-content-between align-items-center">
                <span><?php echo $Lang['first_order_price']; ?></span>
                <?php if($Detail['host_data']['status'] == 'Paid'): if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): ?>
                <button type="button" class="btn btn-primary btn-sm waves-effect waves-light" id="renew" onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['immediate_renewal']; ?></button>
                <?php endif; ?>
                <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                  <button type="button" class="btn btn-primary btn-sm waves-effect waves-light" id="renewpay"><?php echo $Lang['immediate_renewal']; ?></button>
                </a>
                <?php endif; ?>
              </div>
            </div> -->
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray pre-money-box">
                <p class="text-gray">
                  <?php echo $Lang['first_order_price']; ?>
                  <span class="fr"><?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] ==
                    'onetime'): else: ?><?php echo $Detail['host_data']['format_nextduedate']['msg']; ?><?php endif; ?></span>
                </p>
                <section class="d-flex align-items-center">
                  <h3 class="mb-0 mr-2 dc">
                    <?php echo !empty($Detail['host_data']['firstpaymentamount_desc']) ? $Detail['host_data']['firstpaymentamount_desc'] : '-'; ?></h3>

                  <!-- <span class="badge
                      <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?> - <?php else: ?><?php echo $Detail['host_data']['format_nextduedate']['msg']; ?><?php endif; ?></span> -->
                </section>

                <section class="d-flex align-items-center flex-wrap">
                  <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  $Detail['host_data']['status'] == 'Paid'): ?>
                  <span><?php echo $Lang['automatic_balance_renewal']; ?></span>
                  <div class="custom-control custom-switch custom-switch-md mb-4 ml-2" dir="ltr">
                    <input type="checkbox" class="custom-control-input" id="automaticRenewal"
                      onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew'] !=0): ?>checked
                      <?php endif; ?>> <label class="custom-control-label" for="automaticRenewal"></label>
                  </div>
                  <?php endif; if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): if($Detail['host_data']['status'] == 'Paid'): ?>
                  <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb" id="renew"
                    onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                  <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                  <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                    <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb"
                      id="renewpay"><?php echo $Lang['renew']; ?></button>
                  </a>
                  <?php endif; ?>
                  <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): if($Cancel['host_cancel']): ?>
                  <!-- <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn"
                      onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
                      ==
                      'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button> -->
                  <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn" data-container="body"
                    data-toggle="popover" data-placement="top" data-trigger="hover"
                    data-content="将于{<?php echo date('Y-m-d',!is_numeric($Detail['host_data']['deletedate'])? strtotime($Detail['host_data']['deletedate']) : $Detail['host_data']['deletedate']); ?>}自动删除"
                    onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['cancel_out']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-danger btn-sm rsb ml-2" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['apply_out']; ?></button>
                  <?php endif; ?>
                  <?php endif; ?>
                </section>

                <section class="text-gray">
                  <p><?php echo $Lang['subscription_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></p>
                  <p><?php echo $Lang['payment_cycle']; ?>：<?php echo $Detail['host_data']['billingcycle_desc']; ?></p>

                  <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                  <p><?php echo $Lang['due_date']; ?>：-</p>
                  <?php else: ?>
                  <p><?php echo $Lang['due_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></p>
                  <?php endif; ?>
                </section>
              </div>
            </div>
          </div>
          <!-- 订购价格end -->
          <!-- 配置项 -->
          <div class="row">
            <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '5'||$item['option_type'] == '6'||$item['option_type'] == '8'): else: ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['name']; ?>：<?php echo $item['sub_name']; ?>" class="bg-light card-body bg-gray mg-0 row">
                <p class="text-gray col-md-6 plr-0 text-nowrap">
                  <?php echo $item['name']; ?>
                </p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php if($item['option_type']===12): ?>
                  <img src="/upload/common/country/<?php echo $item['code']; ?>.png" width="20px">
                  <?php endif; ?>
                  <?php echo $item['sub_name']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>
            <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['fieldname']; ?>：<?php echo $item['value']; ?>" class="bg-light card-body bg-gray mg-0 row">
                <p class="text-gray col-md-6 plr-0 text-nowrap"><?php echo $item['fieldname']; ?></p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php echo $item['value']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>
            <?php endforeach; ?>
            <div onclick="isShowConfiguration()" class="configuration-btn-down isClick">查看更多信息</div>
            <div class="col-12 mb-2">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $Lang['remarks_infors']; ?>：<?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>"
                class="bg-light card-body  bg-gray mg-0 row">
                <p class="text-gray col-md-3 plr-0 text-nowrap"><?php echo $Lang['remarks_infors']; ?></p>
                <p class="mb-0 col-md-9 plr-0 text-nowrap"><?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>

              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="col-md-9">
      <div class="card">
        <div class="card-body ">
          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">
            <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on'): ?>
            <li class="nav-item" id="chartLi">
              <a class="nav-link active" data-toggle="tab" href="#home1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="fas fa-home"></i></span> -->
                <span><?php echo $Lang['charts']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!($Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on')): ?>active<?php endif; ?>" data-toggle="tab"
                href="#profile1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="far fa-user"></i></span> -->
                <span><?php echo $Lang['upgrade_downgrade']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['download_data']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
              $Detail['host_data']['allow_upgrade_product']))): ?>active<?php endif; ?>" data-toggle="tab" href="#messages1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->
                <span><?php echo $Lang['file_download']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['show_traffic_usage']): ?>
            <li class="nav-item" id="usedLi">
              <a class="nav-link <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
              $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data']): ?>active<?php endif; ?>" data-toggle="tab"
                href="#dosage" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->
                <span><?php echo $Lang['consumption']; ?></span>
              </a>
            </li>
            <?php endif; ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
              $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#settings1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->
                <span><?php echo $Lang['journal']; ?></span>
              </a>
            </li>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#finance" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->
                <span><?php echo $Lang['finance']; ?></span>
              </a>
            </li>
          </ul>

          <!-- Tab panes -->
          <div class="tab-content p-3 text-muted">
            <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on'): ?>
            <div class="tab-pane active" id="home1" role="tabpanel">
            </div>
            <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])): ?>
            <div class="tab-pane <?php if(!($Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on')): ?>active<?php endif; ?>" id="profile1"
              role="tabpanel">
              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-sm-3">

                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>
                        <div class="col-sm-6">

                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>
                        </div>
                        <div class="col-sm-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-sm-3">
                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                        </div>
                        <div class="col-sm-6">
                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>
                        </div>
                        <div class="col-sm-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; if($Detail['download_data']): ?>
            <div class="tab-pane <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product']))): ?>active<?php endif; ?>" id="messages1" role="tabpanel">
              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
            </div>
            <?php endif; ?>
            <div
              class="tab-pane <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
              id="settings1" role="tabpanel">

            </div>

            <?php if($Detail['host_data']['show_traffic_usage']): ?>
            <div class="tab-pane <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data']): ?>active<?php endif; ?>" id="dosage"
              role="tabpanel">
              <div class="row d-flex align-items-center">
                <div class="col-md-3">
                  <input class="form-control" type="date" id="startingTime">
                </div>
                <span><?php echo $Lang['reach']; ?></span>
                <div class="col-md-3">
                  <input class="form-control" type="date" id="endTime">
                </div>
              </div>
              <div class="w-100 h-100">
                <div style="height: 500px" class="chart_content_box w-100" id="usedChartBox"></div>
              </div>
            </div>
            <?php endif; ?>
            <div class="tab-pane" id="finance" role="tabpanel">
              <!-- 财务 -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- 破解密码弹窗 -->
<div class="modal fade" id="dcimModuleResetPass" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0"><?php echo $Lang['crack_password']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <form id="crackPsdForm">
        <div class="modal-body">
          <div class="form-group row mb-0">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['password']; ?></label>
            <div class="col-md-6">
              <input type="text" class="form-control getCrackPsd" name="password">
            </div>
            <div class="col-md-1 fs-18 d-flex align-items-center">
              <i class="fas fa-dice create_random_pass" onclick="create_random_pass()"></i>
            </div>
          </div>
          <label id="password-error-tip" class="control-label error-tip" for="password"></label>
          <div class="form-group row mb-4">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['crack_other']; ?></label>
            <div class="col-md-6">
              <div class="custom-control custom-checkbox">
                <input type="checkbox" class="custom-control-input" id="dcimModuleResetPassOther" value="1"
                  onchange="dcimModuleResetPassOther($(this))">
                <label class="custom-control-label" for="dcimModuleResetPassOther"><?php echo $Lang['password_will_cracked']; ?></label>
              </div>
            </div>
          </div>
          <div class="form-group row mb-4" style="display:none;" id="dcimModuleResetPassOtherUser">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['custom_user']; ?></label>
            <div class="col-md-6">
              <input type="text" class="form-control" name="user">
            </div>
          </div>
        </div>
        <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
      </form>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
        <button type="button" class="btn btn-primary submit"
          onclick="dcimModuleResetPass($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['determine']; ?></button>
      </div>
    </div>
  </div>
</div>

<!-- 救援系统弹窗 -->
<div class="modal fade" id="dcimModuleRescue" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0"><?php echo $Lang['rescue_system']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <form>
        <div class="modal-body">
          <div class="form-group row mb-4">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['system']; ?></label>
            <div class="col-md-8">
              <select class="form-control" name="system">
                <option value="1">Linux</option>
                <option value="2">Windows</option>
              </select>
            </div>
          </div>
        </div>
        <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
      </form>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
        <button type="button" class="btn btn-primary submit"
          onclick="dcimModuleRescue($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['determine']; ?></button>
      </div>
    </div>
  </div>
</div>

<!-- 重装系统弹窗 -->
<div class="modal fade" id="dcimModuleReinstall" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0"><?php echo $Lang['reinstall_system']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <form id="rebuildPsdForm">
        <div class="modal-body">
          <div class="row">
            <div class="col-md-2  d-flex align-items-center justify-content-end">
              <label class="float-right mb-0"><?php echo $Lang['system']; ?></label>
            </div>
            <div class="col-md-5">
              <div class="form-group mb-0">
                <select class="form-control configoption_os_group selectpicker" data-style="btn-default" name="os_group"
                  onchange="dcimModuleReinstallOsGroup($(this))">
                  <?php foreach($Detail['cloud_os_group'] as $item): if(strtolower($item['name'])=="windows"): $os_svg = '1'; elseif(strtolower($item['name'])=="centos"): $os_svg = '2'; elseif(strtolower($item['name'])=="ubuntu"): $os_svg = '3'; elseif(strtolower($item['name'])=="debian"): $os_svg = '4'; elseif(strtolower($item['name'])=="esxi"): $os_svg = '5'; elseif(strtolower($item['name'])=="xenserver"): $os_svg = '6'; elseif(strtolower($item['name'])=="freebsd"): $os_svg = '7'; elseif(strtolower($item['name'])=="fedora"): $os_svg = '8'; else: $os_svg = '9'; ?>
                  <?php endif; ?>
                  <option
                    data-content="<img class='mr-1' src='/upload/common/system/<?php echo $os_svg; ?>.svg' height='20'/><?php echo $item['name']; ?>"
                    value="<?php echo $item['id']; ?>"><?php echo $item['name']; ?></option>
                  <?php endforeach; ?>
                </select>
              </div>
            </div>
            <div class="col-md-5">
              <div class="form-group">
                <select class="form-control" name="os" data-os='<?php echo json_encode($Detail['cloud_os']); ?>'>
                  <?php foreach($Detail['cloud_os'] as $item): ?>
                  <option value="<?php echo $item['id']; ?>" data-group="<?php echo $item['group']; ?>"><?php echo $item['name']; ?></option>
                  <?php endforeach; ?>
                </select>
              </div>
            </div>
          </div>
          <div class="form-group row mb-0">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['password']; ?></label>
            <div class="col-md-6">
              <input type="text" class="form-control getRebuildPsd" name="password">
            </div>
            <div class="col-md-1 fs-18 d-flex align-items-center">
              <i class="fas fa-dice create_random_pass" onclick="create_random_pass()"></i>
            </div>
          </div>
          <label id="password-error-tip-rebuild" class="control-label error-tip ml-8" for="password"></label>
          <div class="row mt-3">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['port']; ?></label>
            <div class="col-md-5">
              <input type="text" class="form-control" name="port" value="22">
            </div>
            <div class="col-md-1 fs-18 d-flex align-items-center">
              <i class="fas fa-dice module_reinstall_random_port"
                onclick="$('#dcimModuleReinstall input[name=\'port\']').val(parseInt(Math.random() * 65535))"></i>
            </div>
          </div>
          <div class="row" id="dcimModuleReinstallPart">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['partition_type']; ?></label>
            <div class="col-md-3">
              <div class="custom-control custom-radio">
                <input type="radio" class="custom-control-input" id="dcimModuleReinstallPart0" name="part_type"
                  onchange="showPartTypeConfirm('<?php echo $Detail['host_data']['os']; ?>')" value="0" checked="checked">
                <label class="custom-control-label" for="dcimModuleReinstallPart0"><?php echo $Lang['full_format']; ?></label>
              </div>
            </div>
            <div class="col-md-3">
              <div class="custom-control custom-radio">
                <input type="radio" class="custom-control-input" id="dcimModuleReinstallPart1" name="part_type"
                  onchange="showPartTypeConfirm('<?php echo $Detail['host_data']['os']; ?>')" value="1">
                <label class="custom-control-label"
                  for="dcimModuleReinstallPart1"><?php echo $Lang['first_partition_formatting']; ?></label>
              </div>
            </div>
            <div class="col-md-3">
              <div class="custom-control custom-checkbox">
                <input type="checkbox" class="custom-control-input" id="dcimModuleReinstallHigh"
                  onchange="showDcimDisk()">
                <label class="custom-control-label" for="dcimModuleReinstallHigh"><?php echo $Lang['senior']; ?></label>
              </div>
            </div>
          </div>
          <div class="row" id="dcimModuleReinstallDisk" style="display:none;">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['disk']; ?></label>
            <div class="col-md-6">
              <select class="form-control" name="disk">
                <option value="0">disk0</option>
              </select>
            </div>
          </div>
          <div class="row">
            <div class="col-md-2 offset-md-2 d-flex align-items-center justify-content-end">
              <div class="custom-control custom-checkbox mb-4">
                <input type="checkbox" class="custom-control-input" id="dcimModuleReinstallConfirm" value="1"
                  onchange="dcimReinstallConfirm($(this))">
                <label class="custom-control-label" for="dcimModuleReinstallConfirm"><?php echo $Lang['finished_backup']; ?></label>
              </div>
            </div>
          </div>
          <div class="row" id="dcimModuleReinstallPartMsg" style="display:none;">
            <div class="col-md-2"></div>
            <div class="col-md-9">
              <div class="part_error"></div>
            </div>
          </div>
          <div class="row">
            <div class="col-md-2"></div>
            <div class="col-md-9">
              <div id="dcimModuleReinstallMsg"></div>
            </div>
          </div>
        </div>
        <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
      </form>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
        <button type="button" class="btn btn-primary submit disabled" style="cursor:not-allowed;"
          onclick="dcimModuleReinstall($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['determine']; ?></button>
      </div>
    </div>
  </div>
</div>
<!-- 电源状态 -->
<script>
  var showPowerStatus = '<?php echo $Detail['module_power_status']; ?>';
  var powerStatus = {}

  $(function () {
    if (showPowerStatus == '1') {
      dcimGetPowerStatus('<?php echo app('request')->get('id'); ?>')
    }

    $('#powerBox').on('click', function () {
      dcimGetPowerStatus('<?php echo app('request')->get('id'); ?>')
    });
  })



  var timeOut = null
  var timeInterval = null
</script>

<script>
  var showPWd = false
  $(function () {
    $(".nav-tabs-custom").find('.nav-item').find("a").removeClass('active')
    $(".nav-tabs-custom").find('.nav-item').eq(0).find("a").addClass('active')
    // 暂停状态悬浮原因
    $('.container-fluid').on('mouseover', '.status-suspended', function () {
      $('.xf-bg').show();
    })
    $('.container-fluid').on('mouseout', '.status-suspended', function () {
      $('.xf-bg').hide();
    })
    // 查看密码
    // $('#copyPwdContent').hide()

    // 复制IP
    var clipboard = new ClipboardJS('.btn-copy-ip', {
      text: function (trigger) {
        return $('#copyIPContent').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    // 一个ip时，不弹框，复制ip
    var clipboard = new ClipboardJS('.copyOneIp', {
      text: function (trigger) {
        return $('#copyOneIp').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    // 复制密码
    var clipboard = new ClipboardJS('.btn-copy-pwd', {
      text: function (trigger) {
        return $('#copyPwdContent').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })

  function togglePwd() {
    showPWd = !showPWd

    if (showPWd) {
      $('#copyPwdContent').show()
      $('#hidePwdBox').hide()
    }
    if (!showPWd) {
      $('#copyPwdContent').hide()
      $('#hidePwdBox').show()
    }
  }
  //点击更多信息
  function isShowConfiguration(first = true) {
    let time = 0
    if (first) {
      time = 500
    }
    for (let i = 0; i < 99; i++) {
      if (i > 3) {
        $('.configuration').eq(i).toggle(time)
      }
    }
    $('.configuration-btn-down').toggleClass('isClick')
    $('.configuration-btn-down').html('查看更多信息')
    $('.configuration-btn-down.isClick').html('收起更多信息')
    if($('.configuration_list').length < 4) {
      $('.configuration-btn-down').hide()
    }
  }
  isShowConfiguration(false)
</script>

<script>
  // 图表tabs
  $(document).ready(function () {
    getComponentData()
  });

  let switch_id = []
  let chartsData = []
  let timeArray = []
  let name = []
  let typeArray = []
  let myChart = null

  $('#chartLi').on('click', function () {
    setTimeout(function () {
      myChart.resize()
    }, 0);
  });
  async function getChartDataFn(index) {
    selectTimeTypeFunc(index)
    const queryObj = {
      id: '<?php echo app('request')->get('id'); ?>',
      switch_id: switch_id[index],
      port_name: name[index],
      start_time: moment(timeArray[index].startTime).valueOf()
    }
    $.ajax({
      type: "post",
      url: '' + '/dcim/traffic',
      data: queryObj,
      success: function (data) {

        var obj = data.data.traffic || []
        var inArray = []
        var outArray = []
        var xName = []
        var inVal = []
        var outVal = []
        for (const item of obj) {
          if (item.type === 'in') {
            inArray.push(item)
            xName.push(moment(item.time).format('MM-DD HH:mm:ss'))
            inVal.push(item.value)
          } else if (item.type === 'out') {
            outArray.push(item)
            outVal.push(item.value)
          }
        }
        chartFunc(index, xName, inVal, outVal, data.data.unit)
        if (data.status === 200) {
          var obj = data.data.traffic
          var inArray = []
          var outArray = []
          var xName = []
          var inVal = []
          var outVal = []
          for (const item of obj) {
            if (item.type === 'in') {
              inArray.push(item)
              xName.push(moment(item.time).format('MM-DD HH:mm:ss'))
              inVal.push(item.value)
            } else if (item.type === 'out') {
              outArray.push(item)
              outVal.push(item.value)
            }
          }
          chartFunc(index, xName, inVal, outVal, data.data.unit)
        }
      }
    });
  }

  async function getComponentData() {
    const obj = {
      id: "<?php echo app('request')->get('id'); ?>"
    }
    $.ajax({
      type: "GET",
      url: '' + '/dcim/detail',
      data: obj,
      success: function (data) {
        if (data.status !== 200) {
          return
        }
        chartsData = data.data.switch ? data.data.switch : []
        let str = ``
        for (let i = 0; i < chartsData.length; i++) {
          const item = chartsData[i];
          timeArray.push({
            startTime: '',
            endTime: ''
          })
          typeArray.push({
            type: '7'
          })
          switch_id.push(chartsData[i].switch_id)
          name.push(chartsData[i].name)
          str += `<div
                    class="w-100 h-100">
                    <select class="form-control" id="chartSelect${i}" class="second_type" name="type" onchange="getChartDataFn(${i})">
                      <option value="24">24<?php echo $Lang['hour']; ?></option>
                      <option value="3">3<?php echo $Lang['day']; ?></option>
                      <option value="7" selected>7<?php echo $Lang['day']; ?></option>
                      <option value="30">30<?php echo $Lang['day']; ?></option>
                      <option value="999"><?php echo $Lang['whole']; ?></option>
                    </select>
                    <div style="height: 500px" class="w-100 h-100 d-flex justify-content-center" id="balanceCharts${i}"></div>
                    </div>`
        }
        $('#home1').append(str);
        for (let j = 0; j < chartsData.length; j++) {
          getChartDataFn(j)
        }
      }
    });
  }

  function selectTimeTypeFunc(index) {
    typeArray[index].type = $(`#chartSelect${index}`).val();

    if (typeArray[index].type === '7') { // 7天
      timeArray[index].startTime = moment(Date.now() - 7 * 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '3') { // 3天
      timeArray[index].startTime = moment(Date.now() - 3 * 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '30') { // 30天
      timeArray[index].startTime = moment(Date.now() - 30 * 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '24') { // 24h
      timeArray[index].startTime = moment(Date.now() - 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '999') { // q全部
      timeArray[index].startTime = ''
    }
  }

  function chartFunc(index, xNameArray, inValArray, outValArray, unitY) {


    var inflow = '<?php echo $Lang['inflow_flow']; ?>'
    var outflow = '<?php echo $Lang['outflow_flow']; ?>'
    // 基于准备好的dom，初始化echarts实例
    myChart = echarts.init(document.getElementById('balanceCharts' + index))
    var option = {
      tooltip: {
        show: true,
        backgroundColor: '#fff',
        borderColor: '#eee',
        showContent: true,
        extraCssText: 'box-shadow: 0 1px 9px rgba(0, 0, 0, 0.1);',
        textStyle: {
          color: '#1e1e2d',
          textBorderWidth: 1
        },
        trigger: 'axis',
        axisPointer: {
          color: '#D9DAEA'
        }
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '10%',
        top: '8%',
        containLabel: true
      },
      color: ['#007bfc', '#3fbf70'],
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: {
        splitLine: {
          show: false
        },
        axisLine: {
          lineStyle: {
            color: '#1e1e2d'
          }
        },
        axisTick: {
          show: false
        },
        type: 'category',
        boundaryGap: false,
        data: xNameArray
      },
      yAxis: {
        axisLabel: {
          // formatter: '{value}' + unitY
          formatter: function (value) {
            return value + unitY
          }
        },
        axisLine: {
          show: false
        },
        minorTick: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#f5f4f8'
          }
        },
        type: 'value'
      },
      series: [{
          name: inflow,
          type: 'line',
          smooth: true,
          data: inValArray
        },
        {
          name: outflow,
          type: 'line',
          smooth: true,
          data: outValArray
        }
      ]
    }
    myChart.setOption(option, true) // true重绘
    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }
</script>

<script>
  // 用量tabs
  let usedChart = null
  let usedStartTime
  let usedEndTime

  $(document).ready(function () {
    if ('<?php echo $Detail['host_data']['show_traffic_usage']; ?>') {
      chartOption()
    }
    if ($('#startingTime,#endTime').length > 0) getData()
    window.addEventListener('resize', function () {
      usedChart.resize()
    })
  })
  $('#usedLi').on('click', function () {
    setTimeout(function () {
      usedChart.resize()
    }, 0);
  });

  $('#startingTime,#endTime').on('change', function () {
    usedStartTime = $('#startingTime').val()
    usedEndTime = $('#endTime').val()
    getData()
  });
  // 图表配置
  function chartOption() {
    usedChart = echarts.init(document.getElementById('usedChartBox'))
    usedChart.setOption({
      backgroundColor: '#fff',
      title: {
        subtext: '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 400
        }
      },
      tooltip: {
        backgroundColor: '#fff',
        padding: [10, 20, 10, 8],
        textStyle: {
          color: '#000',
          fontSize: 12
        },
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        formatter: function (params, ticket, callback) {
          // 
          const res = `
                    <div>
                      <div>` + '<?php echo $Lang['traffic_usage']; ?>' + `：${params[0].value}GB </div>
                      <div>${params[0].axisValue}</div>
                    </div>
                    `
          return res
        },
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
      },
      grid: {
        left: '80',
        top: 20,
        x: 70,
        x2: 50,
        y2: 80
      },
      xAxis: {
        offset: 15,
        type: 'category',
        data: [],
        boundaryGap: false,
        axisTick: {
          show: false
        },
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            type: 'dashed',
            color: '#ddd',
            width: 1
          }
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          }
        }
      },
      yAxis: {
        type: 'value',
        // 轴网格
        splitLine: {
          show: true,
          lineStyle: {
            color: '#ddd',
            type: 'dashed'
          }
        },
        axisTick: {
          show: false // 轴刻度不显示
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          },
          formatter: '{value}GB'
        }
      },
      series: [{
        name: '用量',
        type: 'line',
        smooth: true,
        showSymbol: true,
        symbol: 'circle',
        symbolSize: 3,
        // data: ['1200', '1400', '1008', '1411', '1026', '1288', '1300', '800', '1100', '1000', '1118', '123456'],
        data: [],
        areaStyle: {
          normal: {
            color: '#d4d1da',
            opacity: 0.2
          }
        },
        itemStyle: {
          normal: {
            color: '#0061ff' // 主要线条的颜色
          }
        },
        lineStyle: {
          normal: {
            width: 4,
            shadowColor: 'rgba(0,0,0,0.4)',
            shadowBlur: 10,
            shadowOffsetY: 10
          }
        }
      }]
    })
  }

  // 获取数据
  async function getData() {
    usedChart.showLoading({
      text: '<?php echo $Lang['data_loading']; ?>' + '...',
      color: '#999',
      textStyle: {
        fontSize: 30,
        color: '#444'
      },
      effectOption: {
        backgroundColor: 'rgba(0, 0, 0, 0)'
      }
    })



    const obj = {
      id: '<?php echo app('request')->get('id'); ?>',
      start: usedStartTime,
      end: usedEndTime
    }
    $.ajax({
      type: "get",
      url: '' + '/dcim/traffic_usage',
      data: obj,
      success: function (data) {

        usedChart.hideLoading()
        if (data.status !== 200) return false
        const xData = []
        const seriesData = [];
        (data.data || []).forEach(item => {
          xData.push(item.time)
          seriesData.push(item.value)
        })
        usedChart.setOption({
          title: {
            subtext: xData.length ? '' : '<?php echo $Lang['no_data_available']; ?>'
          },
          xAxis: {
            data: xData
          },
          series: [{
            data: seriesData
          }]
        })
        // 如果初始查询没有时间, 则设置默认时间为返回数据的第一个和最后一个时间
        if (!usedStartTime || !usedEndTime) {
          if (data.data.length) {
            usedStartTime = new Date().getFullYear() + '-' + data.data[0].time
            usedEndTime = new Date().getFullYear() + '-' + data.data[data.data.length - 1].time
            $('#startingTime').val(new Date().getFullYear() + '-' + data.data[0].time);
            $('#endTime').val(new Date().getFullYear() + '-' + data.data[data.data.length - 1].time);
          }
        }
      }
    });
  }

  // 分辨率改变, 重绘图表
  function resize() {
    usedChart.resize()
  }

  // 时间选择改变
  function dateChange() {
    const startTimeStamp = new Date(usedStartTime).getTime()
    const endTimeStamp = new Date(usedEndTime).getTime()
    if (usedStartTime && usedEndTime && startTimeStamp < endTimeStamp) {
      getData()
    }
  }
</script>

<script>
  // 获取基础数据
  const obj = {
    host_id: '<?php echo app('request')->get('id'); ?>'
  }
  $.ajax({
    type: "get",
    url: '' + '/host/dedicatedserver',
    data: obj,
    success: function (data) {
      const totalFlow = data.data.host_data.bwlimit // 总流量
      const usedFlow = data.data.host_data.bwusage // 已用流量
      const remainingFlow = (totalFlow - usedFlow).toFixed(1)
      let percentUsed = 100 - parseInt((usedFlow / totalFlow) * 100) || 0
      $('#totalProgress')
        .css('width', percentUsed + '%')
        .attr('aria-valuenow', percentUsed)
        .text(`${percentUsed}%`);

      $('#usedFlowSpan').text(`${usedFlow > 1024 ? ((usedFlow / 1024).toFixed(2) + 'TB') : (usedFlow + 'GB')}`);
      $('#remainingFlow').text(
        `${remainingFlow > 1024 ? ((remainingFlow / 1024).toFixed(2) + 'TB') : (remainingFlow + 'GB')}`);

      // 产品状态
      $('#statusBox').append(`<span class="sprite2 ${data.data.host_data.domainstatus}"></span>`)
    }
  });
</script>

<script>
  const logObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'log_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: logObj,
    success: function (data) {
      $(data).appendTo('#settings1');
    }
  });
  // 财务的append
  const financeObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'billing_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: financeObj,
    success: function (data) {
      $(data).appendTo('#finance');
    }
  });
</script>

<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<script>
  var getResintallStatusTimer = null;
  $(function () {
    getResintallStatus('<?php echo app('request')->get('id'); ?>');
  })
</script>

<script>
  var clipboard = null
  var clipboardpoppwd = null
  var ips = <?php echo json_encode($Detail['host_data']['assignedips']); ?>;
  // console.log('ips: ', ips);
  $(document).on('click', '#copyIPContent', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('IP地址')
    var iplist = ''
    if (clipboard) {
      clipboard.destroy()
    }
    for(let item in ips) {
      iplist += `
        <div>
          <span class="copyIPContent${item}">${ips[item]}</span>
          <i class="bx bx-copy pointer text-primary ml-1 btn-copy btnCopyIP${item}" data-clipboard-action="copy" data-clipboard-target=".copyIPContent${item}"></i>
        </div>
      `

      // 复制IP
      clipboard = new ClipboardJS('.btnCopyIP'+item, {
        text: function (trigger) {
          return $('.copyIPContent'+item).text()
        },
        container: document.getElementById('popModal')
      });
      clipboard.on('success', function (e) {
        toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
      })
    }

    $('#popContent').html(iplist)
  });


  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      
    `)
  });


  $('#popModal').on('shown.bs.modal', function () {
    if (clipboardpoppwd) {
      clipboardpoppwd.destroy()
    }
    clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })

  // 破解密码
  $(document).on('blur', '.getCrackPsd', function () {
    veriCrackPsd()
  })

  function veriCrackPsd() {
    let result = checkingPwd($(".getCrackPsd").val(), passwordRules.num, passwordRules.upper, passwordRules.lower,
      passwordRules.special)
    if (result.flag) {
      $('#password-error-tip').css('display', 'none');
      $('.getCrackPsd').removeClass("is-invalid");
    } else {
      $("#password-error-tip").html(result.msg);
      $(".getCrackPsd").addClass("is-invalid");
      $('#password-error-tip').css('display', 'block');
    }
  }
  //重装系统
  $(document).on('blur', '.getRebuildPsd', function () {
    veriRebuildPsd()
  })

  function veriRebuildPsd() {
    let result = checkingPwd($(".getRebuildPsd").val(), passwordRules.num, passwordRules.upper, passwordRules.lower,
      passwordRules.special)
    if (result.flag) {
      $('#password-error-tip-rebuild').css('display', 'none');
      $('.getRebuildPsd').removeClass("is-invalid");
    } else {
      $("#password-error-tip-rebuild").html(result.msg);
      $(".getRebuildPsd").addClass("is-invalid");
      $('#password-error-tip-rebuild').css('display', 'block');
    }
  }

  $(function () {
    $("#crackPsdForm").on('click', ".create_random_pass", function (e) {
      veriCrackPsd()
    })
    $("#rebuildPsdForm").on('click', ".create_random_pass", function (e) {
      veriRebuildPsd()
    })
  })
</script>
<?php elseif($Detail['host_data']['type'] == "cloud"): ?>
    <style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<script src="themes/clientarea/default/assets/libs/echarts/echarts.min.js?v=<?php echo $Ver; ?>"></script>
<style>
  .server_header_box {
    height: auto;
    background-image: linear-gradient(87deg, #4d83ff 0%, #3656ff 100%);
    border-radius: 15px;
    padding: 20px 25px;
    color: #ffffff;
  }

  .left_wrap_btn {
    display: inline-block;
    width: 80px;
    height: 20px;
    background-color: #5f88fe;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    color: #ffffff;
    text-align: center;
    border: none;
  }

  .custom-button {
    background-color: #6f87fc;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    font-size: 12px;
    color: #fff;
    border: none;
  }

  .box_left_wrap {
    border-left: 1px solid rgba(255, 255, 255, 0.25);
    min-height: 74px;
  }

  @media screen and (max-width: 1367px) {
    .form-control {
      width: 46%;
    }

    .server_header_box {
      height: auto;
    }

    .power_box {
      max-width: 300px;
    }

    .left_wrap_btn {
      width: 60px !important;
    }

    .bottom-box {
      margin-top: 3rem !important;
    }

    .osbox {
      max-width: 150px;
    }
  }

  @media screen and (max-width: 976px) {
    .server_header_box {
      height: auto;
      padding: 20px;
      margin-top: 10px;
    }

    .domain,
    .box_left_wrap {
      margin-bottom: 20px;
      border-left: none;
    }

    .power_box {
      margin-bottom: 20px;
    }
  }

  .tuxian {
    cursor: pointer;
  }

  .tuxian:hover {
    color: rgba(224, 224, 224, 0.877);
  }

  #copyIPContent:hover {
    color: #FCA426
  }

  #copyOneIp:hover {
    color: #FCA426
  }

  .dc {
    color: #5f88fe
  }

  .rsb {
    height: 20px;
    padding: 0px 10px;
  }

  .mg-0 {
    margin: 0;
  }

  .plr-0 {
    padding-left: 0px;
    padding-right: 0px;
    margin-bottom: 0;
  }

  .text-nowrap {
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .text-right {
    text-align: right;
  }

  .pre-money-box {
    background: url("/themes/clientarea/default/assets/images/money.png") no-repeat;
    background-position-x: right;
    background-position-y: bottom;
  }

  .w-75 {
    width: 75% !important;
  }

  .ll-flex {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .xf-bg {
    display: none;
    position: absolute;
    background: #fff;
    padding: 10px 15px;
    border-radius: 4px;
    top: -40px;
    left: 0px;
    box-shadow: 0px 3px 5px 0px rgba(0, 28, 144, 0.21);
    font-size: 12px;
  }

  .xf-bg-text {
    color: #333;
    word-break: break-all;
  }

  .restall-btn {
    border-radius: 25px;
    margin-left: 20px;
  }

  .fr {
    float: right;
  }

  .flex-wrap {
    display: flex !important;
    flex-flow: wrap;
  }
  .configuration-btn-down {
  width: 100%;

  text-align: center;
  line-height: 36px;
  color: #5F88FE;
  }

  .configuration-btn-down::after {
  content: '';
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-left: 8px;
  background-color: transparent;
  transform: rotate(225deg);
  border: 1px solid #5F88FE;
  border-bottom: none;
  border-right: none;
  transform-origin: 2px;
  transition: all .2s;
  }

  .configuration-btn-down.isClick::after {
  transform: rotate(45deg);
  }
</style>
<div class="container-fluid">

  <div class="row mb-4">

    <div class="col-12">

      <div class="row align-items-center server_header_box">

        <div class="mr-3 power_box">

          <div class="text-white d-flex">
            <!-- 电源状态 -->
            <div class="mr-3 pointer">
              <?php if($Detail['module_power_status'] == '1'): ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="powerBox">
                <span id="powerStatusIcon" class="bx bx-loader" data-toggle="popover" data-trigger="hover" title=""
                  data-html="true" data-content="<?php echo $Lang['loading']; ?>..."></span>
              </div>
              <?php else: ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="statusBox"></div>
              <?php endif; ?>
            </div>
            <div>
              <section class="d-flex align-items-center mb-2">

                <h4 class="text-white mb-0 font-weight-bold"><?php echo $Detail['host_data']['productname']; ?></h4>

                <span class="badge badge-pill ml-2 py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>">

                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>

                </span>

              </section>

              <section>

                <span><?php echo $Detail['host_data']['domain']; ?></span>

              </section>
            </div>

          </div>

        </div>

        <div class="pl-4 mr-3 box_left_wrap osbox">

          <span class="text-white-50 fs-12"><?php echo $Lang['operating_system']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $Detail['host_data']['os']; ?></h5>

          <?php if(in_array('reinstall', array_column($Detail['module_button']['control'], 'func')) &&
          $Detail['host_data']['domainstatus']=="Active"): ?>

          <span class="ml-0">

            <button type="button" class="service_module_button left_wrap_btn fs-12 restall-btn" data-func="reinstall"
              data-type="default"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reinstall_system']; ?></button>

          </span>

          <?php endif; ?>

        </div>

        <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '6'||$item['option_type'] == '8'): ?>

        <div class="pl-4 mr-3 box_left_wrap">

          <span class="text-white-50 fs-12"><?php echo $item['name']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $item['sub_name']; ?></h5>

        </div>

        <?php endif; ?>

        <?php endforeach; ?>

        <div class="pl-4 mr-3 box_left_wrap">
          <span class="text-white-50 fs-12"><?php echo $Lang['ip_address']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white">
            <!-- <span data-toggle="popover" data-trigger="hover" title="" data-html="true" data-content="
                  <?php foreach($Detail['host_data']['assignedips'] as $list): ?>
                  <div><?php echo $list; ?></div>
                  <?php endforeach; ?>
                "> -->
            <span>
              <?php if($Detail['host_data']['dedicatedip']): if($Detail['host_data']['assignedips']): ?>
              <span class="tuxian" id="copyIPContent"
                class="pointer"><?php echo $Detail['host_data']['dedicatedip']; ?>(<?php echo count($Detail['host_data']['assignedips']); ?>)</span>
              <!-- <span>(<?php echo count($Detail['host_data']['assignedips']); ?>)</span> -->
              <?php else: ?>
              <span id="copyOneIp" class="pointer copyOneIp"><?php echo $Detail['host_data']['dedicatedip']; ?></span>
              <?php endif; ?>
              <!-- <i class="bx bx-copy pointer text-white ml-1 btn-copy" id="btnCopyIP" data-clipboard-action="copy"
                  data-clipboard-target="#copyIPContent"></i> -->
              <?php else: ?>
              -
              <?php endif; ?>

            </span>

          </h5>



        </div>

        <!--
        <div class="pl-4 mr-3 box_left_wrap">

          <span class="text-white-50 fs-12"><?php echo $Lang['password']; ?></span>
          <h5 class="mt-2" data-toggle="popover" data-trigger="hover" data-html="true"
            data-content="<?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?><br><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?>">
            <span id="hidePwdBox" class="text-white">***********</span>
            <span id="copyPwdContent" class="text-white"><?php echo $Detail['host_data']['password']; ?></span>
            <i class="fas fa-eye pointer ml-2 text-white" onclick="togglePwd()"></i>
            <i class="bx bx-copy pointer ml-1 btn-copy text-white" id="btnCopyPwd" data-clipboard-action="copy"
              data-clipboard-target="#copyPwdContent"></i>
          </h5>

          

          <?php if(in_array('crack_pass', array_column($Detail['module_button']['control'], 'func')) && $Detail['host_data']['domainstatus']=="Active"): ?>

          <button type="button" class="service_module_button ml-0 left_wrap_btn fs-12" data-func="crack_pass"
            data-type="default"
            onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reset_password']; ?></button>

          <?php endif; ?>

        </div>
        -->

        <div class="d-flex justify-content-end flex-shrink-1 flex-grow-1">

          <?php if($Detail['module_button']['control'] && $Detail['host_data']['domainstatus']=="Active"): ?>

          <div class="btn-group ml-2 mr-2 mt-2">

            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['control'] as $item): if($item['func'] != 'crack_pass' && $item['func'] != 'reinstall'): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

              <?php endif; ?>

              <?php endforeach; ?>

            </div>

          </div>

          <?php endif; if($Detail['module_button']['console'] && $Detail['host_data']['domainstatus']=="Active"): ?>

          <div class="btn-group ml-2 mr-2 mt-2">
            <?php if(($Detail['module_button']['console']|count) == 1): foreach($Detail['module_button']['console'] as $item): ?>
            <a class="btn btn-primary service_module_button d-flex align-items-center custom-button"
              href="javascript:void(0);"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
              data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
            <?php endforeach; else: ?>
            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['console'] as $item): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

              <?php endforeach; ?>

            </div>
            <?php endif; ?>

          </div>

          <?php endif; ?>

          <!-- <div class="btn-group ml-2 mr-2 mt-2">

            <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>

            <span>

              <?php if($Cancel['host_cancel']): ?>

              <button class="btn btn-danger mb-1 h-100" id="cancelStopBtn"
                onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
                ==
                'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button>

              <?php else: ?>

              <button class="btn btn-primary mb-1 h-100 custom-button" data-toggle="modal"
                data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>

              <?php endif; ?>

            </span>

            <?php endif; ?>
          </div> -->

          <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.2-->
          <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
          <div class="btn-group ml-2 mr-2 mt-2">
            <span>
              <?php echo $item; ?>
            </span>
          </div>
          <?php endforeach; ?>
          <?php endif; ?>
          <!-- 结束 -->
        </div>

      </div>

    </div>

  </div>

  <div class="row bottom-box">

    <div class="col-md-3">

      <div class="card">

        <div class="card-body">

          <!-- <div class="mb-3 text-center" data-toggle="popover" data-trigger="hover" data-html="true"
            data-placement="bottom"
            data-content="<?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?><br><?php echo $Lang['password']; ?>：<?php echo $Detail['host_data']['password']; ?><br><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?>"
            style="width: 100px;height: 30px;line-height: 30px;background-color: #ffffff;box-shadow: 0px 4px 20px 2px rgba(6, 75, 179, 0.08);border-radius: 4px;cursor: pointer;">
          <div class="mb-3 text-center" id="logininfo" style="width: 100px;height: 30px;line-height: 30px;background-color: #ffffff;box-shadow: 0px 4px 20px 2px rgba(6, 75, 179, 0.08);border-radius: 4px;cursor: pointer;">
            <?php echo $Lang['login_information']; ?>
          </div>
          -->
          <!-- 登录信息 -->
          <?php if($Detail['host_data']['domainstatus'] == 'Active'): ?>
          <div class="row">
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['login_information']; ?>
                  <i class="bx bx-user dc"></i>
                </p>
                <p class="mb-0"><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></p>
                <p class="mb-0"><?php echo $Lang['password']; ?>：
                  <!-- <span id="hidePwdBox" class="text-black">***********</span> -->
                  <?php if($Detail['host_data']['password'] == ''): ?>
                  <span class="text-black btnCopyPwd pointer dc">-</span>
                  <?php else: ?>
                  <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制"
                    id="copyPwdContent" class="text-black btnCopyPwd pointer dc"><?php echo $Detail['host_data']['password']; ?></span>
                  <?php endif; ?>
                  <!-- <i class="fas fa-eye pointer ml-2 text-black" onclick="togglePwd()"></i> -->
                  <!-- <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPwd" data-clipboard-action="copy"
                    data-clipboard-target="#copyPwdContent"></i> -->

                  <?php if(in_array('crack_pass', array_column($Detail['module_button']['control'], 'func')) &&
                  $Detail['host_data']['domainstatus']=="Active"): ?>

                  <button type="button" class="service_module_button ml-0 left_wrap_btn fs-12 fr" data-func="crack_pass"
                    data-type="default" 
                    onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reset_password']; ?></button>

                  <?php endif; ?>
                </p>
                <p class="mb-0"><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] ==
                  '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></p>

              </div>
            </div>
          </div>
          <?php endif; ?>

          <!-- 流量start -->
          <?php if(($Detail['host_data']['domainstatus'] == 'Active' || ($Detail['host_data']['domainstatus'] == 'Suspended' && $Detail['host_data']['suspendreason_type'] == 'flow')) && $Detail['host_data']['bwlimit'] > 0): ?>
          <!-- <div class="d-flex justify-content-end mb-2">

            <button type="button" class="btn btn-success btn-sm waves-effect waves-light" id="orderFlowBtn"
              onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>

          </div> -->
          <div class="mt-4 mb-3">
            <i class="bx bx-circle" style="color:#f0ad4e"></i> <?php echo $Lang['used_flow']; ?>：<span id="usedFlowSpan">-</span>
            <i class="bx bx-circle" style="color:#34c38f"></i> <?php echo $Lang['residual_flow']; ?>：<span id="remainingFlow">-</span>
          </div>
          <div class="mb-4 ll-flex">
            <div class="progress w-75">
              <div class="progress-bar progress-bar-striped bg-success" id="totalProgress" role="progressbar"
                style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">100%
              </div>
            </div>
            <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb ml-2" id="orderFlowBtn"
              onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
          </div>
          <?php endif; ?>
          <!-- 流量end -->
          <!-- 订购价格 -->
          <div class="row">
            <!-- <div class="col-12 my-2">
              <div class="d-flex justify-content-between align-items-center">
                <span><?php echo $Lang['first_order_price']; ?></span>
                <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): if($Detail['host_data']['status'] == 'Paid'): ?>
                <button type="button" class="btn btn-primary btn-sm waves-effect waves-light" id="renew"
                  onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['immediate_renewal']; ?></button>
                <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                  <button type="button" class="btn btn-primary btn-sm waves-effect waves-light"
                    id="renewpay"><?php echo $Lang['immediate_renewal']; ?></button>
                </a>
                <?php endif; ?>
                <?php endif; ?>
              </div>
            </div> -->

            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray pre-money-box">
                <p class="text-gray">
                  <?php echo $Lang['first_order_price']; ?>
                  <span class="fr"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): else: ?>
                    <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span>
                </p>
                <section class="d-flex align-items-center">
                  <h3 class="mb-0 mr-2 dc">
                    <?php echo !empty($Detail['host_data']['firstpaymentamount_desc']) ? $Detail['host_data']['firstpaymentamount_desc'] : '-'; ?></h3>
                  <!-- <span class="badge
                      <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): ?> - <?php else: ?> <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span> -->
                </section>

                <?php if($Renew['host']['billingcycle'] != 'onetime' && $Renew['host']['status'] == 'Paid' && $Renew['host']['billingcycle'] !=
                'free'): ?>
                <section class="d-flex align-items-center flex-wrap">
                  <span><?php echo $Lang['automatic_balance_renewal']; ?></span>
                  <div class="custom-control custom-switch custom-switch-md mb-4 ml-2" dir="ltr">
                    <input type="checkbox" class="custom-control-input" id="automaticRenewal"
                      onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew'] !=0): ?>checked
                      <?php endif; ?>> <label class="custom-control-label" for="automaticRenewal"></label>
                  </div>
                  <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): if($Detail['host_data']['status'] == 'Paid'): ?>
                  <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb" id="renew"
                    onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                  <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                  <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                    <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb"
                      id="renewpay"><?php echo $Lang['renew']; ?></button>
                  </a>
                  <?php endif; ?>
                  <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>
                  <span>
                    <?php if($Cancel['host_cancel']): ?>
                    <!-- <button class="btn btn-danger mb-1 h-100" id="cancelStopBtn"
                          onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
                          ==
                          'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button> -->
                    <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn" data-container="body"
                      data-toggle="popover" data-placement="top" data-trigger="hover"
                      data-content="将于{<?php echo date('Y-m-d',!is_numeric($Detail['host_data']['deletedate'])? strtotime($Detail['host_data']['deletedate']) : $Detail['host_data']['deletedate']); ?>}自动删除"
                      onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['cancel_out']; ?></button>
                    <?php else: ?>
                    <button class="btn btn-danger btn-sm rsb ml-2" data-toggle="modal"
                      data-target=".cancelrequire"><?php echo $Lang['apply_out']; ?></button>
                    <?php endif; ?>
                  </span>
                  <?php endif; ?>
                </section>
                <?php endif; ?>


                <section class="text-gray">

                  <p><?php echo $Lang['subscription_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></p>

                  <p><?php echo $Lang['payment_cycle']; ?>：<?php echo $Detail['host_data']['billingcycle_desc']; ?></p>

                  <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                  <p><?php echo $Lang['due_date']; ?>：<?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): ?> - <?php else: ?>
                    <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></p>
                  <?php else: ?>
                  <p><?php echo $Lang['due_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></p>
                  <?php endif; ?>

                </section>

              </div>

            </div>
          </div>


          <!-- 远程地址 建站解析 -->
          <?php if($Detail['dcimcloud'] && ($Detail['dcimcloud']['nat_acl'] || $Detail['dcimcloud']['nat_web'])): ?>
          <div class="bg-primary rounded-sm d-flex flex-column justify-content-center text-white py-2 px-1 mb-2">
            <?php if($Detail['dcimcloud']['nat_acl']): ?>
            <div class="d-flex justify-content-between align-items-center">
              <span>
                <label><?php echo $Lang['remote_address']; ?>：</label>
                <span id="nat_aclBox"><?php echo $Detail['dcimcloud']['nat_acl']; ?></span>
              </span>
              <span>
                <i class="bx bx-copy pointer text-white btn-copy" id="btnCopyaclBox" data-clipboard-action="copy"
                  data-clipboard-target="#nat_aclBox"></i>
              </span>
            </div>
            <?php endif; if($Detail['dcimcloud']['nat_web']): ?>
            <div class="d-flex justify-content-between align-items-center">
              <span>
                <label><?php echo $Lang['ip_address']; ?>建站解析：</label>
                <span id="nat_webBox"><?php echo $Detail['dcimcloud']['nat_web']; ?></span>
              </span>
              <span>
                <i class="bx bx-copy pointer text-white btn-copy" id="btnCopywebBox" data-clipboard-action="copy"
                  data-clipboard-target="#nat_webBox"></i>
              </span>
            </div>
            <?php endif; ?>
          </div>
          <?php endif; ?>
          <!-- 配置项 -->
          <div class="row">
            <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '5'||$item['option_type'] == '6'||$item['option_type'] == '8'): else: ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['name']; ?>：<?php echo $item['sub_name']; ?>" class="bg-light card-body py-2 bg-gray row">
                <p class="text-gray col-md-6 plr-0 text-nowrap">
                  <?php echo $item['name']; ?>
                </p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php if($item['option_type']===12): ?>
                  <img src="/upload/common/country/<?php echo $item['code']; ?>.png" width="20px">
                  <?php endif; ?>
                  <?php echo $item['sub_name']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>

            <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['fieldname']; ?>：<?php echo $item['value']; ?>" class="bg-light card-body py-2 bg-gray row">
                <p class="text-gray col-md-6 plr-0 text-nowrap"><?php echo $item['fieldname']; ?></p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php echo $item['value']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>
            <?php endforeach; ?>
            <div onclick="isShowConfiguration()" class="configuration-btn-down isClick">查看更多信息</div>
            <div class="col-12 mb-2">

              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $Lang['remarks_infors']; ?>：<?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>"
                class="bg-light card-body bg-gray row">

                <p class="text-gray col-md-3 plr-0 text-nowrap"><?php echo $Lang['remarks_infors']; ?></p>

                <p class="mb-0 col-md-9 plr-0 text-nowrap"><?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>

              </div>

            </div>
          </div>

        </div>

      </div>

    </div>

    <div class="col-md-9">

      <div class="card">

        <div class="card-body">

          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">

            <?php if($Detail['module_chart']): ?>

            <li class="nav-item" id="chartLi">

              <a class="nav-link active" data-toggle="tab" href="#home1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-home"></i></span> -->

                <span><?php echo $Lang['charts']; ?></span>

              </a>

            </li>

            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>

            <li class="nav-item">

              <a class="nav-link <?php if(!$Detail['module_chart']): ?>active<?php endif; ?>" data-toggle="tab" href="#profile1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="far fa-user"></i></span> -->

                <span><?php echo $Lang['upgrade_downgrade']; ?></span>

              </a>

            </li>

            <?php endif; if($Detail['download_data']): ?>

            <li class="nav-item">

              <a class="nav-link <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#messages1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->

                <span><?php echo $Lang['file_download']; ?></span>

              </a>

            </li>

            <?php endif; if($Detail['host_data']['show_traffic_usage']): ?>
            <li class="nav-item" id="usedLi">

              <a class="nav-link <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#dosage" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->

                <span><?php echo $Lang['consumption']; ?></span>

              </a>

            </li>
            <?php endif; foreach($Detail['module_client_area'] as $key=>$item): ?>

            <li class="nav-item">

              <a class="nav-link  <?php if($key==0 && !$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data']): ?>active<?php endif; ?>" data-toggle="tab" href="#module_client_area_<?php echo $item['key']; ?>" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

                <span><?php echo $item['name']; ?></span>

              </a>

            </li>

            <?php endforeach; ?>

            <li class="nav-item">

              <a class="nav-link <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#settings1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

                <span><?php echo $Lang['journal']; ?></span>

              </a>

            </li>

            <li class="nav-item">

              <a class="nav-link" data-toggle="tab" href="#finance" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

                <span><?php echo $Lang['finance']; ?></span>

              </a>

            </li>

          </ul>



          <!-- Tab panes -->

          <div class="tab-content p-3 text-muted" style="min-height: 550px;">

            <?php if($Detail['module_chart']): ?>

            <div class="tab-pane active" id="home1" role="tabpanel">

              <style>
    .chartbox {
      width: 50%!important;
    }
    @media screen and (max-width: 1367px) {
      .chartbox {
        width: 100%!important;
      }
    }
</style>
<div class="row">
<?php foreach($Detail['module_chart'] as $key=>$item): ?>

  <div class="chartbox">
    <div class="module_chart module_chart_<?php echo $item['type']; ?>" data-type="<?php echo $item['type']; ?>" data-title="<?php echo $item['title']; ?>">
      <div style="height: 60px;">
        <h4 style="float: left;"><?php echo $item['title']; ?></h4>
        <?php if($item['select']): ?>
        <div style="width:200px;float: right;margin-right: 20px;">
          <select class="module_chart_select_<?php echo $item['type']; ?> form-control selectpicker_refresh" onchange="getChartDataFn('','<?php echo $key; ?>')">
            <?php foreach($item['select'] as $item2): ?>
            <option value="<?php echo $item2['value']; ?>"><?php echo $item2['name']; ?></option>
            <?php endforeach; ?>
          </select>
        </div>
        <?php endif; ?>
      </div>
      <div class="module_chart_date">
        <input class="form-control startTime" type="datetime-local" onchange="getChartDataFn('','<?php echo $key; ?>')">
        <span class="ml-1 mr-1"><?php echo $Lang['reach']; ?></span>
        <input class="form-control endTime mr-3" type="datetime-local" onchange="getChartDataFn('','<?php echo $key; ?>')">
      </div>

    </div>

    <div class="w-100 h-100 ">

      <div style="height: 500px" class="chart_content_box w-100" id="module_chart_<?php echo $item['type']; ?>"></div>

    </div>
  </div>
<?php endforeach; ?>

</div>

  




<script>
  // 图表tabs
  $(document).ready(function () {
    var arr = JSON.parse('<?php echo json_encode($Detail['module_chart']); ?>')
    setTimeout(function(){
      arr.forEach(function(item){
        getChartDataFn(item)
      })
    }, 0);
    
  });

  let switch_id = []
  let chartsData = []
  let timeArray = []
  let name = []
  let typeArray = []
  let myChart = null

  $('#chartLi').on('click', function () {
    setTimeout(function(){
      myChart.resize()
    }, 0);
  });


  // line
  function lineChartOption (type, xAxisData, seriesData0, seriesData1, unit, label) {
    // 硬盘IO
    const myChart = echarts.init(document.getElementById('module_chart_'+type))
    myChart.setOption({
      backgroundColor: '#fff',
      title: {
        subtext: (!xAxisData.length) ? '暂无数据' : '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 250
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        backgroundColor: '#fff',
        textStyle: {
          color: '#333',
          fontSize: 12
        },
        padding: [10, 10],
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);',
        formatter: function (params, ticket, callback) {
          // console.log('line:', params)
          const res = `<div>
                        <div>${params[0].marker} ${params[0].seriesName}：${params[0].value}${unit}</div>
                        <div>${params[1] ? params[1].marker : ''} ${params[1] ? params[1].seriesName : ''}${params[1] ? '：' : ''}${params[1] ? params[1].value : ''}${params[1] ? unit : ''}</div>
                        <div style="color: #999999;">${params[0].axisValue}</div>
                      </div>`
          return res
        }
      },
      grid: {
        left: '80',
        top: 30,
        x: 50,
        x2: 50,
        y2: 80
      },
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: [{
        offset: 15,
        type: 'category',
        boundaryGap: false,
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            type: 'dashed',
            color: '#ddd',
            width: 1
          }
        },
        // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
        //   return str.replace(' ', '\n')
        // }),
        data: xAxisData,
        // 轴刻度
        axisTick: {
          show: false
        },
        // 轴网格
        splitLine: {
          show: false
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999999'
          }
        }
      }],
      yAxis: [{
        type: 'value',
        axisTick: {
          show: false
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          formatter: '{value}' + unit,
          textStyle: {
            color: '#556677'
          }
        },
        splitLine: {
          lineStyle: {
            type: 'dashed'
          }
        }
      }],
      series: [{
        name: label[0],
        type: 'line',
        // data: [5, 12, 11, 14, 25, 16, 10, 18, 6],
        data: seriesData0,
        symbolSize: 1,
        symbol: 'circle',
        smooth: true,
        showSymbol: false,
        lineStyle: {
          width: 2
        },
        itemStyle: {
          normal: {
            color: '#75db16',
            borderColor: '#75db16'
          }
        }
      }, {
        name: label[1],
        type: 'line',
        // data: [10, 10, 30, 12, 15, 3, 7, 20, 15000],
        data: seriesData1,
        symbolSize: 1,
        symbol: 'circle',
        smooth: true,
        showSymbol: false,
        lineStyle: {
          width: 3,
          shadowColor: 'rgba(92, 102, 255, 0.3)',
          shadowBlur: 10,
          shadowOffsetY: 20
        },
        itemStyle: {
          normal: {
            color: '#5c66ff',
            borderColor: '#5c66ff'
          }
        }
      }
      ]
    })

    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }
  // area
  function areaChartOption (type, xAxisData, seriesData, unit, label) {
    // CPU使用率
    const myChart = echarts.init(document.getElementById('module_chart_'+type ))
    myChart.setOption({
      grid: {
        left: '80',
        top: 30,
        x: 50,
        x2: 50,
        y2: 80
      },
      backgroundColor: '#fff',
      title: {
        subtext: (!xAxisData.length) ? '暂无数据' : '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 250
        }
      },
      tooltip: {
        backgroundColor: '#fff',
        padding: [10, 20, 10, 8],
        textStyle: {
          color: '#333',
          fontSize: 12
        },
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        formatter: function (params, ticket, callback) {
          // console.log(params, '')
          const res = `<div>
                        <div>${params[0].seriesName}：${params[0].value}${unit} </div>
                        <div style="color: #999999;">${params[0].axisValue}</div>
                      </div>`
          return res
        },
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
      },
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: {
        offset: 15,
        type: 'category',
        boundaryGap: false,
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            color: '#999999',
            width: 1
          }
        },
        // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
        //   return str.replace(' ', '\n')
        // }),
        data: xAxisData,
        // 轴刻度
        axisTick: {
          show: false
        },
        // 轴网格
        splitLine: {
          show: false
        },
        axisLabel: {
          show: true,
          // interval: 0, // 横轴信息全部显示
          textStyle: {
            color: '#999999'
          }
        }
      },
      yAxis: {
        axisTick: {
          show: false // 轴刻度不显示
        },
        max: 100,
        min: 0,
        // 改变y轴颜色
        axisLine: {
          show: false
        },
        // 轴网格
        splitLine: {
          show: true,
          lineStyle: {
            color: '#ddd',
            type: 'dashed'
          }
        },
        // 坐标轴文字样式
        axisLabel: {
          show: true,
          formatter: '{value}' + unit,
          textStyle: {
            color: '#999999'
          }
        }
      },
      series: [{
        name: label,
        type: 'line',
        areaStyle: {
          opacity: 1,
          color: '#737dff'
        },
        symbol: 'none', // 折线无拐点
        lineStyle: {
          normal: {
            width: 0 // 折线宽度
          }
        },
        smooth: true,
        // data: [5, 25, 20, 50, 10, 40, 18, 25, 0]
        data: seriesData

      }]
    })

    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }

  // bar
  function barChartOption (type, xAxisData, seriesData0, seriesData1, unit, label) {
    // 内存用量
    const myChart = echarts.init(document.getElementById('module_chart_'+type))
    myChart.setOption({
      backgroundColor: '#fff',
      title: {
        subtext: (!xAxisData.length) ? '暂无数据' : '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 250
        }
      },
      tooltip: {
        backgroundColor: '#fff',
        padding: [10, 20, 10, 8],
        textStyle: {
          color: '#000',
          fontSize: 12
        },
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        formatter: function (params, ticket, callback) {
          // console.log('bar:', params)
          const res = `
          <div>
              <div>${params[0].marker}${params[0].seriesName}：${params[0].value}${unit} </div>                
              <div>${params[1] ? params[1].marker : ''} ${params[1] ? params[1].seriesName : ''}${params[1] ? '：' : ''}${params[1] ? params[1].value : ''}${params[1] ? unit : ''}</div>
              <div>${params[0].axisValue}</div>
          </div>`
          return res
        },
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
      },
      grid: {
        left: '80',
        top: 30,
        x: 70,
        x2: 50,
        y2: 80
      },
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: {
        offset: 15,
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          }
        },
        type: 'category',
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            type: 'dashed',
            color: '#ddd',
            width: 1
          }
        },
        // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
        //   return str.replace(' ', '\n')
        // })
        data: xAxisData
      },
      yAxis: {
        axisTick: {
          show: false // 轴刻度不显示
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          },
          formatter: '{value}' + unit
        },
        // 轴网格
        splitLine: {
          show: true,
          lineStyle: {
            color: '#ddd',
            type: 'dashed'
          }
        }

      },
      series: [{
        name: label[1],
        type: 'bar',
        stack: '总量',
        barGap: '-100%',
        // data: [136, 132, 101, 134, 90, 230, 210, 100, 300],
        data: seriesData1,
        itemStyle: {
          barBorderRadius: [5, 5, 0, 0],
          color: '#737dff'
        }
      },
      {
        name: label[0],
        type: 'bar',
        stack: '可用',
        // data: [964, 182, 191, 234, 290, 330, 310, 100, 500],
        data: seriesData0,
        itemStyle: {
          barBorderRadius: [5, 5, 0, 0],
          color: '#ccc',
          opacity: 0.3
        }
      } 
      ]
    })

    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }



  async function getChartDataFn(e,index) {
    

    var arr = JSON.parse('<?php echo json_encode($Detail['module_chart']); ?>')
    if(index){e = arr[index]}
    var echartTT = echarts.init(document.getElementById('module_chart_'+e.type))
    echartTT.showLoading({
      text: '数据正在加载...',
      color: '#999',
      textStyle: {
        fontSize: 30,
        color: '#444'
      },
      effectOption: {
        backgroundColor: 'rgba(0, 0, 0, 0)'
      }
    })
    const queryObj = {
      id: '<?php echo app('request')->get('id'); ?>',
      type: e.type,
      start: new Date($('.module_chart_'+e.type+' .startTime').val()).getTime(),
      end:new Date($('.module_chart_'+e.type+' .endTime').val()).getTime(),
      select: $('.module_chart_select_' + e.type).val(),
    }
    $.ajax({
      type: "GET",
      url: '' + '/provision/chart/<?php echo app('request')->get('id'); ?>',
      data: queryObj,
      success: function (data) {
        echartTT.hideLoading()
      if (data.status !== 200) return false

      const xAxisData = []
      const seriesData0 = []
      const seriesData1 = [];

      (data.data.list || []).forEach((item, index) => {
        (item || []).forEach(innerItem => {
          if (index === 0) {
            xAxisData.push(innerItem.time)
            seriesData0.push(innerItem.value)
          } else if (index === 1) {
            seriesData1.push(innerItem.value)
          }
        })
      })

       if (data.data.chart_type === 'area') {
          areaChartOption(e.type, xAxisData, seriesData0, data.data.unit, data.data.label)
        } else if (data.data.chart_type === 'line') {
          lineChartOption(e.type, xAxisData, seriesData0, seriesData1, data.data.unit, data.data.label)
        } else if (data.data.chart_type === 'bar') {
          barChartOption(e.type, xAxisData, seriesData0, seriesData1, data.data.unit, data.data.label)
        }

        // 如果初始查询没有时间, 则设置默认时间为返回数据的第一个和最后一个时间
        if (!$('.module_chart_'+e.type+' .startTime').val() || !$('.module_chart_'+e.type+' .endTime').val()) {
          if (data.data.list[0].length) {
            var start = new Date(data.data.list[0][0].time).getTime()
            var end = new Date(data.data.list[0][data.data.list[0].length - 1].time).getTime()
            $('.module_chart_'+e.type+' .startTime').val(moment(start).format('YYYY-MM-DDTHH:mm'))
            $('.module_chart_'+e.type+' .endTime').val(moment(end).format('YYYY-MM-DDTHH:mm'))
          }
        }
      }
    });
  }





</script>

            </div>

            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <div class="tab-pane <?php if(!$Detail['module_chart']): ?>active<?php endif; ?>" id="profile1" role="tabpanel">

              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">

                  <div class="col-12">

                    <div class="bg-light  rounded card-body">

                      <div class="row">

                        <div class="col-sm-3">
                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>

                        <div class="col-sm-6">



                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>

                        </div>

                        <div class="col-sm-3">

                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>

                        </div>

                      </div>

                    </div>

                  </div>

                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">

                  <div class="col-12">

                    <div class="bg-light  rounded card-body">

                      <div class="row">

                        <div class="col-sm-3">

                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>

                        </div>

                        <div class="col-sm-6">

                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>

                        </div>

                        <div class="col-sm-3">

                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')">
                            
                            <?php echo $Lang['upgrade_downgrade_options']; ?>
                          </button>

                        </div>

                      </div>

                    </div>

                  </div>

                </div>
                <?php endif; ?>

              </div>

            </div>
            <?php endif; if($Detail['download_data']): ?>

            <div
              class="tab-pane <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
              id="messages1" role="tabpanel">

              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>

            </div>

            <?php endif; foreach($Detail['module_client_area'] as $key=>$item): ?>

            <div class="tab-pane <?php if($key==0 && !$Detail['module_chart'] && !$Detail['download_data'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>" role="tabpanel" id="module_client_area_<?php echo $item['key']; ?>">
              <div style="min-height: 550px;width:100%">
                <script>
                  $.ajax({
                    url : '/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>'
                    ,type : 'get'
                    ,success : function(res) {
                        $('#module_client_area_<?php echo $item['key']; ?> > div').html(res);
                    }
                  })
                </script>
              </div>
              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                onload="this.height=$($('.main-content .card-body')[1]).height()-72" frameborder="0"
                width="100%"></iframe> -->
              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                frameborder="0" width="100%" style="min-height: 550px;"></iframe> -->

            </div>

            <?php endforeach; if($Detail['host_data']['show_traffic_usage']): ?>
            <div
              class="tab-pane <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data']): ?>active<?php endif; ?>"
              id="dosage" role="tabpanel">

              <div class="row d-flex align-items-center">

                <div class="col-md-3">

                  <input class="form-control" type="date" id="startingTime">

                </div>

                <span><?php echo $Lang['reach']; ?></span>

                <div class="col-md-3">

                  <input class="form-control" type="date" id="endTime">

                </div>

              </div>

              <div class="w-100 h-100">

                <div style="height: 500px" class="chart_content_box w-100" id="usedChartBox"></div>

              </div>

            </div>
            <?php endif; ?>
            <div
              class="tab-pane <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage'] && !$Detail['module_client_area']): ?>active<?php endif; ?>"
              id="settings1" role="tabpanel">

              <!-- 日志 -->

            </div>
            <div class="tab-pane" id="finance" role="tabpanel">

              <!-- 财务 -->

            </div>

          </div>

        </div>

      </div>

    </div>

  </div>

</div>

<!-- 电源状态 -->
<script>
  var showPowerStatus = '<?php echo $Detail['module_power_status']; ?>';
  var powerStatus = {
    status: '',
    des: ''
  }

  $(function () {
    if (showPowerStatus == '1') {
      getPowerStatus('<?php echo app('request')->get('id'); ?>')
    }

    $('#powerBox').on('click', function () {
      getPowerStatus('<?php echo app('request')->get('id'); ?>')
    });
  })

  // 获取电源状态
  function getPowerStatus() {

    $('#powerStatusIcon').removeClass()
    $('#powerStatusIcon').addClass('bx bx-loader')

    $.ajax({
      type: "POST",
      url: '/provision/default',
      data: {
        id: '<?php echo app('request')->get('id'); ?>',
        func: 'status'
      },
      success: function (data) {
        $('#powerStatusIcon').attr('data-content', data.data ? data.data.des : data.msg)
        $('#powerStatusIcon').removeClass()
        if (data.status != 200) {
          powerStatus.status = 'unknown'
          powerStatus.des = data.msg
          $('#powerStatusIcon').addClass('sprite unknown')
        } else {
          powerStatus.status = data.data ? data.data.status : 'unknown'
          powerStatus.des = data.data ? data.data.des : '<?php echo $Lang['unknown']; ?>'
          // 

          if (powerStatus.status === 'process') {
            $('#powerStatusIcon').addClass('bx bx-loader')
          } else if (powerStatus.status === 'on') {
            $('#powerStatusIcon').addClass('sprite start')
          } else if (powerStatus.status === 'off') {
            $('#powerStatusIcon').addClass('sprite closed')
          } else if (powerStatus.status === 'waiting') {
            $('#powerStatusIcon').addClass('sprite waitOn')
          } else if (powerStatus.status === 'suspend') {
            $('#powerStatusIcon').addClass('sprite pause')
          } else if (powerStatus.status === 'wait_reboot' || powerStatus.status === 'wait') {
            $('#powerStatusIcon').addClass('sprite waiting')
          } else if (powerStatus.status === 'cold_migrate') {
            $('#powerStatusIcon').addClass('iconfont icon-shujuqianyi')
          } else if (powerStatus.status === 'hot_migrate') {
            $('#powerStatusIcon').addClass('iconfont icon-shujuqianyi')
          } else {
            $('#powerStatusIcon').addClass('sprite unknown')
          }

          if (data.data.status !== 'process') { // 状态改变, 清除定时器
            clearInterval(timeInterval)
          }
        }

      }
    });
  }

  var timeOut = null
  var timeInterval = null
</script>
<script>
  var showPWd = true
  $(function () {
    $(".nav-tabs-custom").find('.nav-item').find("a").removeClass('active')
    $(".nav-tabs-custom").find('.nav-item').eq(0).find("a").addClass('active')
    // 暂停状态悬浮原因
    $('.container-fluid').on('mouseover', '.status-suspended', function () {
      $('.xf-bg').show();
    })
    $('.container-fluid').on('mouseout', '.status-suspended', function () {
      $('.xf-bg').hide();
    })
    let ips = <?php echo json_encode($Detail['host_data']['assignedips']); ?>;
    if(ips.length<2) {
      // 复制IP
      var clipboard1 = new ClipboardJS('#copyIPContent', {
        text: function (trigger) {
          return $('#copyIPContent').text()
        }
      });
      clipboard1.on('success', function (e) {
        toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
      })
    }
    // 一个ip时，不弹框，复制ip
    var clipboard = new ClipboardJS('.copyOneIp', {
      text: function (trigger) {
        return $('#copyOneIp').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })

    // 查看密码
    // $('#copyPwdContent').hide()
    $('#hidePwdBox').hide()

    // 复制密码
    var clipboard2 = new ClipboardJS('.btnCopyPwd', {
      text: function (trigger) {
        return $('#copyPwdContent').text()
      }
    });
    clipboard2.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })

  function togglePwd() {
    showPWd = !showPWd

    if (showPWd) {
      $('#copyPwdContent').show()
      $('#hidePwdBox').hide()
    }
    if (!showPWd) {
      $('#copyPwdContent').hide()
      $('#hidePwdBox').show()
    }
  }
  //点击更多信息
  function isShowConfiguration(first = true) {
    let time = 0
    if (first) {
      time = 500
    }
    for (let i = 0; i < 99; i++) {
      if (i > 3) {
        $('.configuration').eq(i).toggle(time)
      }
    }
    $('.configuration-btn-down').toggleClass('isClick')
    $('.configuration-btn-down').html('查看更多信息')
    $('.configuration-btn-down.isClick').html('收起更多信息')
    if($('.configuration_list').length < 4) {
      $('.configuration-btn-down').hide()
    }
  }
  isShowConfiguration(false)
  // 远程地址 建站解析
  // 复制授权码
  var clipboardnat_aclBox = new ClipboardJS('#btnCopyaclBox', {
    text: function (trigger) {
      return $('#nat_aclBox').text()
    }
  });
  clipboardnat_aclBox.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

  // 复制授权码
  var clipboardnat_webBox = new ClipboardJS('#btnCopywebBox', {
    text: function (trigger) {
      return $('#nat_webBox').text()
    }
  });
  clipboardnat_webBox.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })
</script>


<script>
  // 用量tabs

  let usedChart = null

  let usedStartTime

  let usedEndTime



  $(document).ready(function () {
    chartOption()

    if ($('#startingTime,#endTime').length > 0) getData()

    window.addEventListener('resize', function () {

      usedChart.resize()

    })



    $('#usedLi').on('click', function () {

      setTimeout(() => {

        usedChart.resize()

      }, 0);

    });


    $('#startingTime,#endTime').change(function () {

      usedStartTime = $('#startingTime').val()

      usedEndTime = $('#endTime').val()

      getData()

    });


    // 图表配置

    function chartOption() {
      if ($('#usedChartBox').length == 0) {
        return;
      }

      usedChart = echarts.init(document.getElementById('usedChartBox'))

      usedChart.setOption({

        backgroundColor: '#fff',

        title: {

          subtext: '',

          left: 'center',

          textAlign: 'left',

          subtextStyle: {

            lineHeight: 400

          }

        },

        tooltip: {

          backgroundColor: '#fff',

          padding: [10, 20, 10, 8],

          textStyle: {

            color: '#000',

            fontSize: 12

          },

          trigger: 'axis',

          axisPointer: {

            type: 'line',

            lineStyle: {

              color: '#7dcb8f'

            }

          },

          formatter: function (params, ticket, callback) {

            // 

            const res = `

                    <div>

                      <div><?php echo $Lang['traffic_usage']; ?>：${params[0].value}GB </div>

                      <div>${params[0].axisValue}</div>

                    </div>

                    `

            return res

          },

          extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'

        },

        grid: {

          left: '80',

          top: 20,

          x: 70,

          x2: 50,

          y2: 80

        },

        xAxis: {

          offset: 15,

          type: 'category',

          data: [],

          boundaryGap: false,

          axisTick: {

            show: false

          },

          // 改变x轴颜色

          axisLine: {

            lineStyle: {

              type: 'dashed',

              color: '#ddd',

              width: 1

            }

          },

          axisLabel: {

            show: true,

            textStyle: {

              color: '#999'

            }

          }

        },

        yAxis: {

          type: 'value',

          // 轴网格

          splitLine: {

            show: true,

            lineStyle: {

              color: '#ddd',

              type: 'dashed'

            }

          },

          axisTick: {

            show: false // 轴刻度不显示

          },

          axisLine: {

            show: false

          },

          axisLabel: {

            show: true,

            textStyle: {

              color: '#999'

            },

            formatter: '{value}GB'

          }

        },

        series: [{

          name: '用量',

          type: 'line',

          smooth: true,

          showSymbol: true,

          symbol: 'circle',

          symbolSize: 3,

          // data: ['1200', '1400', '1008', '1411', '1026', '1288', '1300', '800', '1100', '1000', '1118', '123456'],

          data: [],

          areaStyle: {

            normal: {

              color: '#d4d1da',

              opacity: 0.2

            }

          },

          itemStyle: {

            normal: {

              color: '#0061ff' // 主要线条的颜色

            }

          },

          lineStyle: {

            normal: {

              width: 4,

              shadowColor: 'rgba(0,0,0,0.4)',

              shadowBlur: 10,

              shadowOffsetY: 10

            }

          }

        }]

      })

    }



    // 获取数据

    async function getData() {

      usedChart.showLoading({

        text: '数据正在加载...',

        color: '#999',

        textStyle: {

          fontSize: 30,

          color: '#444'

        },

        effectOption: {

          backgroundColor: 'rgba(0, 0, 0, 0)'

        }

      })







      const obj = {

        id: '<?php echo app('request')->get('id'); ?>',

        start: usedStartTime,

        end: usedEndTime

      }

      $.ajax({

        type: "get",

        url: '/host/trafficusage',

        data: obj,

        success: function (data) {



          usedChart.hideLoading()

          if (data.status !== 200) return false

          const xData = []

          const seriesData = [];

          (data.data || []).forEach(item => {

            xData.push(item.time)

            seriesData.push(item.value)

          })

          usedChart.setOption({

            title: {

              subtext: xData.length ? '' : '暂无数据'

            },

            xAxis: {

              data: xData

            },

            series: [{

              data: seriesData

            }]

          })

          // 如果初始查询没有时间, 则设置默认时间为返回数据的第一个和最后一个时间

          if (!usedStartTime || !usedEndTime) {

            if (data.data.length) {

              usedStartTime = new Date().getFullYear() + '-' + data.data[0].time

              usedEndTime = new Date().getFullYear() + '-' + data.data[data.data.length - 1].time

              $('#startingTime').val(new Date().getFullYear() + '-' + data.data[0].time);

              $('#endTime').val(new Date().getFullYear() + '-' + data.data[data.data.length - 1].time);

            }

          }

        }

      });

    }



    // 分辨率改变, 重绘图表

    function resize() {

      usedChart.resize()

    }



    // 时间选择改变

    function dateChange() {

      const startTimeStamp = new Date(usedStartTime).getTime()

      const endTimeStamp = new Date(usedEndTime).getTime()

      if (usedStartTime && usedEndTime && startTimeStamp < endTimeStamp) {

        getData()

      }

    }
  })
</script>



<script>
  // 获取基础数据

  const obj = {

    host_id: '<?php echo app('request')->get('id'); ?>'

  }

  $.ajax({

    type: "get",

    url: '/host/dedicatedserver',

    data: obj,

    success: function (data) {

      const totalFlow = data.data.host_data.bwlimit // 总流量

      const usedFlow = data.data.host_data.bwusage.toFixed(1) // 已用流量

      const remainingFlow = (totalFlow - usedFlow).toFixed(1)

      let percentUsed = 100 - parseInt((usedFlow / totalFlow) * 100) || 0

      $('#totalProgress')
        .css('width', percentUsed + '%')
        .attr('aria-valuenow', percentUsed)
        .text(`${percentUsed}%`);

      $('#usedFlowSpan').text(`${usedFlow > 1024 ? ((usedFlow / 1024).toFixed(2) + 'TB') : (usedFlow + 'GB')}`);
      $('#remainingFlow').text(
        `${remainingFlow > 1024 ? ((remainingFlow / 1024).toFixed(2) + 'TB') : (remainingFlow + 'GB')}`);

      // 产品状态
      $('#statusBox').append(`<span class="sprite2 ${data.data.host_data.domainstatus}"></span>`)

    }

  });
</script>



<script>
  const logObj = {

    id: '<?php echo app('request')->get('id'); ?>',

    action: 'log_page'

  }

  $.ajax({

    type: "get",

    url: '/servicedetail',

    data: logObj,

    success: function (data) {

      $(data).appendTo('#settings1');

    }

  });

  // 财务的append

  const financeObj = {

    id: '<?php echo app('request')->get('id'); ?>',

    action: 'billing_page'

  }

  $.ajax({

    type: "get",

    url: '/servicedetail',

    data: financeObj,

    success: function (data) {

      $(data).appendTo('#finance');

    }

  });
</script>


<script>
  var clipboard = null
  var clipboardpoppwd = null
  var ips = <?php echo json_encode($Detail['host_data']['assignedips']); ?>;

  $(document).on('click', '#copyIPContent', function () {
    // if (ips.length<=1) {
    //         return
    //     }
    $('#popModal').modal('show')
    $('#popTitle').text('IP地址')
    var iplist = ''
    if (clipboard) {
      clipboard.destroy()
    }
    for(let item in ips) {
      iplist += `
        <div>
          <span class="copyIPContent${item}">${ips[item]}</span>
          <i class="bx bx-copy pointer text-primary ml-1 btn-copy btnCopyIP${item}" data-clipboard-action="copy" data-clipboard-target=".copyIPContent${item}"></i>
        </div>
      `

      // 复制IP
      clipboard = new ClipboardJS('.btnCopyIP'+item, {
        text: function (trigger) {
          return $('.copyIPContent'+item).text()
        },
        container: document.getElementById('popModal')
      });
      clipboard.on('success', function (e) {
        toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
      })
    }

    $('#popContent').html(iplist)
  });


  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      <div>
        <?php if(in_array('crack_pass', array_column($Detail['module_button']['control'], 'func')) && $Detail['host_data']['domainstatus']=="Active"): ?>
        <button type="button" class="service_module_button ml-2 mt-2 left_wrap_btn fs-12" data-func="crack_pass"
          data-type="default"
          onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reset_password']; ?></button>
        <?php endif; ?>
      </div>
    `)
  });


  $('#popModal').on('shown.bs.modal', function () {
    if (clipboardpoppwd) {
      clipboardpoppwd.destroy()
    }
    clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })
</script>

<?php elseif($Detail['host_data']['type'] == "dcimcloud"): ?>
    
<style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<style>
  .server_header_box {
    height: auto;
    background-image: linear-gradient(87deg, #4d83ff 0%, #3656ff 100%);
    border-radius: 15px;
    padding: 20px 25px;
    color: #ffffff;
  }

  .left_wrap_btn {
    display: inline-block;
    width: 80px;
    height: 20px;
    background-color: #5f88fe;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    color: #ffffff;
    text-align: center;
    border: none;
  }

  .custom-button {
    background-color: #6f87fc;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    font-size: 12px;
    color: #fff;
    border: none;
  }

  .box_left_wrap {
    border-left: 1px solid rgba(255, 255, 255, 0.25);
    min-height: 74px;
  }

  .aibiao a {
    width: 100%;
    height: 100%;
    display: inline-block;
  }

  @media screen and (max-width: 1367px) {
    .form-control {
      width: 46%;
    }

    .server_header_box {
      height: auto;
    }

    .power_box {
      max-width: 300px;
    }

    .left_wrap_btn {
      width: 60px !important;
    }

    .bottom-box {
      margin-top: 3rem !important;
    }

    .osbox {
      max-width: 150px;
    }
  }

  @media screen and (max-width: 976px) {
    .server_header_box {
      height: auto;
      padding: 20px;
      margin-top: 10px;
    }

    .domain,
    .box_left_wrap {
      margin-bottom: 20px;
      border-left: none;
    }

    .power_box {
      margin-bottom: 20px;
    }
  }

  .tuxian {
    cursor: pointer;
  }

  .tuxian:hover {
    color: rgba(224, 224, 224, 0.877);
  }

  .alarm {
    display: inline-block;
    font-size: 12px;
    cursor: pointer;
    color: #e3ecf1;
    font-weight: 300;
  }

  .fr {
    float: right;
  }

  .restall-btn {
    border-radius: 25px;
    margin-left: 20px;
  }

  .login-info-icon {
    color: #5f88fe;
  }

  .dc {
    color: #5f88fe
  }

  .rsb {
    height: 20px;
    padding: 0px 10px;
  }

  .mg-0 {
    margin: 0;
  }

  .plr-0 {
    padding-left: 0px;
    padding-right: 0px;
    margin-bottom: 0;
  }

  #copyIPContent:hover {
    color: #FCA426
  }

  #copyOneIp:hover {
    color: #FCA426
  }

  .text-nowrap {
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .text-right {
    text-align: right;
  }

  .pre-money-box {
    background: url("/themes/clientarea/default/assets/images/money.png") no-repeat;
    background-position-x: right;
    background-position-y: bottom;
  }

  .w-75 {
    width: 75% !important;
  }

  .ll-flex {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .xf-bg {
    display: none;
    position: absolute;
    background: #fff;
    padding: 10px 15px;
    border-radius: 4px;
    top: -40px;
    left: 0px;
    box-shadow: 0px 3px 5px 0px rgba(0, 28, 144, 0.21);
    font-size: 12px;
  }

  .xf-bg-text {
    color: #333;
    word-break: break-all;
  }

  .sj {
    display: block;
    width: 10px;
    height: 10px;
    background: #fff;
    position: absolute;
    bottom: -3px;
    transform: rotate(45deg);
  }

  .flex-wrap {
    display: flex !important;
    flex-flow: wrap;
  }

  .configuration-btn-down {
    width: 100%;

    text-align: center;
    line-height: 36px;
    color: #5F88FE;
  }

  .configuration-btn-down::after {
    content: '';
    display: inline-block;
    width: 8px;
    height: 8px;
    margin-left: 8px;
    background-color: transparent;
    transform: rotate(225deg);
    border: 1px solid #5F88FE;
    border-bottom: none;
    border-right: none;
    transform-origin: 2px;
    transition: all .2s;
  }

  .configuration-btn-down.isClick::after {
    transform: rotate(45deg);
  }
  .custom-control-input:disabled+.custom-control-label {
    cursor: not-allowed;
  }

</style>
<!--  -->
<div class="container-fluid">

  <div class="row mb-4">

    <div class="col-12">

      <div class="row align-items-center server_header_box">

        <div class="mr-3 power_box">

          <div class="text-white d-flex">
            <!-- 电源状态 -->
            <div class="mr-3 pointer">
              <?php if($Detail['module_power_status'] == '1'): ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="powerBox">
                <span id="powerStatusIcon" class="bx bx-loader" data-toggle="popover" data-trigger="hover" title=""
                  data-html="true" data-content="<?php echo $Lang['loading']; ?>..."></span>
              </div>
              <?php else: ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="statusBox"></div>
              <?php endif; ?>
            </div>
            <div>
              <section class="d-flex align-items-center mb-2">

                <h4 class="text-white mb-0 font-weight-bold"><?php echo $Detail['host_data']['productname']; ?></h4>

                <span class="badge badge-pill ml-2 py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>"
                  style="position: relative;">

                  <div class="xf-bg">
                    <div class="xf-bg-text"><span style="color: #e31519;">暂停原因：</span><?php echo !empty($Detail['host_data']['suspendreason']) ? $Detail['host_data']['suspendreason'] : $Detail['host_data']['suspendreason_type']; ?></div>
                    <font class="sj"></font>
                  </div>
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>

                </span>

              </section>

              <section>

                <span><?php echo $Detail['host_data']['domain']; ?></span>

              </section>
            </div>

          </div>

        </div>

        <div class="pl-4 mr-3 box_left_wrap osbox">
          <span class="text-white-50 fs-12"><?php echo $Lang['operating_system']; ?></span>
          <?php if(in_array('reinstall', array_column($Detail['module_button']['control'], 'func')) &&
          $Detail['host_data']['domainstatus']=="Active"): ?>
          <span class="ml-0">
            <button type="button" class="service_module_button left_wrap_btn fs-12 restall-btn" data-func="reinstall"
              data-type="default"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reinstall_system']; ?></button>
          </span>
          <?php endif; ?>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $Detail['host_data']['os']; ?></h5>
        </div>

        <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '6'||$item['option_type'] == '8'): ?>

        <div class="pl-4 mr-3 box_left_wrap">

          <span class="text-white-50 fs-12"><?php echo $item['name']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $item['sub_name']; ?></h5>



        </div>

        <?php endif; ?>

        <?php endforeach; ?>

        <div class="pl-4 mr-3 box_left_wrap">

          <span class="text-white-50 fs-12"><?php echo $Lang['ip_address']; ?></span>

          <h5 class="mt-2 font-weight-bold text-white">
            <!-- <span data-toggle="popover" data-trigger="hover" title="" data-html="true" data-content="
                  <?php foreach($Detail['host_data']['assignedips'] as $list): ?>
                  <div><?php echo $list; ?></div>
                  <?php endforeach; ?>
                "> -->
            <span>
              <?php if($Detail['host_data']['dedicatedip']): if($Detail['host_data']['assignedips']): ?>
              <!-- <span class="tuxian"><?php echo $Detail['host_data']['dedicatedip']; ?></span>-->
              <span id="copyIPContent" class="pointer"><?php echo $Detail['host_data']['dedicatedip']; ?>
                (<?php echo count($Detail['host_data']['assignedips']); ?>)</span>
              <?php else: ?>
              <span id="copyOneIp" class="pointer copyOneIp"><?php echo $Detail['host_data']['dedicatedip']; ?></span>
              <?php endif; ?>
              <!-- <?php if(count($Detail['host_data']['assignedips']) >= 0): ?>
                <i class="bx bx-copy pointer text-white ml-1 btn-copy" id="btnCopyIP" data-clipboard-action="copy"
                  data-clipboard-target="#copyIPOne"></i>
                <?php endif; ?> -->
              <?php else: ?>
              -
              <?php endif; ?>
            </span>

            <!-- <?php if(count($Detail['host_data']['assignedips']) > 0): ?>
            <div class="alarm" id="copyIPContent">
              更多
            </div>
            <?php endif; ?> -->

          </h5>


        </div>

        <!--
          <div class="pl-4 mr-3 box_left_wrap">
          <span class="text-white-50 fs-12"><?php echo $Lang['password']; ?></span>
          <h5 class="mt-2" data-toggle="popover" data-trigger="hover" data-html="true"
            data-content="<?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?><br><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?>">
            <span id="hidePwdBox" class="text-white">***********</span>
            <span id="copyPwdContent" class="text-white"><?php echo $Detail['host_data']['password']; ?></span>
            <i class="fas fa-eye pointer ml-2 text-white" onclick="togglePwd()"></i>
            <i class="bx bx-copy pointer ml-1 btn-copy text-white" id="btnCopyPwd" data-clipboard-action="copy"
              data-clipboard-target="#copyPwdContent"></i>
          </h5>

        </div>
        -->

        <div class="d-flex justify-content-end flex-shrink-1 flex-grow-1">

          <?php if($Detail['module_button']['control'] && $Detail['host_data']['domainstatus']=="Active"): ?>

          <div class="btn-group ml-2 mr-2 mt-2">

            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['control'] as $item): if($item['func'] != 'crack_pass' && $item['func'] != 'reinstall'): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

              <?php endif; ?>

              <?php endforeach; ?>

            </div>

          </div>

          <?php endif; if($Detail['module_button']['console']): ?>

          <div class="btn-group ml-2 mr-2 mt-2">
            <?php if(count($Detail['module_button']['console']) == 1): foreach($Detail['module_button']['console'] as $item): ?>
            <a class="btn btn-primary service_module_button d-flex align-items-center custom-button"
              href="javascript:void(0);"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
              data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
            <?php endforeach; else: ?>
            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['console'] as $item): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

              <?php endforeach; ?>

            </div>
            <?php endif; ?>

          </div>

          <?php endif; ?>

          <!-- <div class="btn-group ml-2 mr-2 mt-2">

            <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>

            <span>

              <?php if($Cancel['host_cancel']): ?>

              <button class="btn btn-danger mb-1 h-100" id="cancelStopBtn"
                onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
                ==
                'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button>

              <?php else: ?>

              <button class="btn btn-primary mb-1 h-100 custom-button " data-toggle="modal"
                data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>
              <?php endif; ?>

            </span>

            <?php endif; ?>



          </div> -->

          <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.1-->
          <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
          <div class="btn-group ml-2 mr-2 mt-2">
            <!-- <button class="btn btn-primary mb-1 h-100 custom-button"  data-toggle="modal"
                data-target=".pullserver"><?php echo $item; ?></button> -->
            <?php echo $item; ?>
          </div>
          <?php endforeach; ?>
          <?php endif; ?>
          <!-- 结束 -->

        </div>

      </div>

    </div>

  </div>
  <div class="row bottom-box">

    <div class="col-md-3">

      <div class="card">

        <div class="card-body">

          <!-- <div class="mb-3 text-center" data-toggle="popover" data-trigger="hover" data-html="true"
            data-placement="bottom"
            data-content="<?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?><br><?php echo $Lang['password']; ?>：<?php echo $Detail['host_data']['password']; ?><br><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?>"
            style="width: 100px;height: 30px;line-height: 30px;background-color: #ffffff;box-shadow: 0px 4px 20px 2px rgba(6, 75, 179, 0.08);border-radius: 4px;cursor: pointer;">
          <div class="mb-3 text-center" id="logininfo" style="width: 100px;height: 30px;line-height: 30px;background-color: #ffffff;box-shadow: 0px 4px 20px 2px rgba(6, 75, 179, 0.08);border-radius: 4px;cursor: pointer;">
            <?php echo $Lang['login_information']; ?>
          </div> -->

          <div class="row">
            <!-- 登录信息start -->
            <?php if($Detail['host_data']['domainstatus'] == 'Active'): ?>
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['login_information']; ?>
                  <i class="bx bx-user login-info-icon"></i>
                </p>
                <p class="mb-0"><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></p>
                <p class="mb-0"><?php echo $Lang['password']; ?>：
                  <!-- <span id="hidePwdBox" class="text-black">***********</span> -->
                  <?php if($Detail['host_data']['password'] == ''): ?>
                  <span class="text-black btnCopyPwd pointer dc">-</span>
                  <?php else: ?>
                  <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制"
                    id="copyPwdContent" class="text-black btnCopyPwd pointer dc"><?php echo $Detail['host_data']['password']; ?></span>
                  <?php endif; ?>
                  <!-- <i class="fas fa-eye pointer ml-2 text-black" onclick="togglePwd()"></i> -->
                  <!-- <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPwd" data-clipboard-action="copy"
                    data-clipboard-target="#copyPwdContent"></i> -->

                  <?php if(in_array('crack_pass', array_column($Detail['module_button']['control'], 'func')) &&
                  $Detail['host_data']['domainstatus']=="Active"): ?>

                  <button type="button" class="service_module_button ml-0 left_wrap_btn fs-12 fr" data-func="crack_pass"
                    data-type="default" style="white-space: nowrap;"
                    onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reset_password']; ?></button>

                  <?php endif; ?>
                </p>
                <p class="mb-0"><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] ==
                  '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></p>

              </div>
            </div>
            <?php endif; ?>
            <!-- 登录信息end -->
            <!-- 面板管理信息start -->
            <?php if(($temp_custom_field_data = array_column($Detail['custom_field_data'], 'value', 'fieldname')) &&
            (isset($temp_custom_field_data['panel_address']) || isset($temp_custom_field_data['面板管理地址']) ||
            isset($temp_custom_field_data['panel_passwd']) || isset($temp_custom_field_data['面板管理密码']))): ?>
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['panel_manage_info']; ?>
                  <i class="bx bx-receipt dc"></i>
                </p>
                <?php if(isset($temp_custom_field_data['panel_address']) || isset($temp_custom_field_data['面板管理地址'])): ?>
                <p class="mb-0"><?php echo $Lang['panel_manage_address']; ?>：<?php echo isset($temp_custom_field_data['panel_address']) ? $temp_custom_field_data['panel_address'] : $temp_custom_field_data['面板管理地址']; ?></p>
                <?php endif; if(isset($temp_custom_field_data['panel_passwd']) || isset($temp_custom_field_data['面板管理密码'])): ?>
                <p class="mb-0"><?php echo $Lang['panel_manage_password']; ?>：
                  <!-- <span id="hidePanelPasswd">***********</span> -->
                  <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制"
                    id="panelPasswd" class="btnCopyPanelPasswd dc pointer"><?php echo isset($temp_custom_field_data['panel_passwd']) ? $temp_custom_field_data['panel_passwd'] : ($temp_custom_field_data['面板管理密码'] ?:'--'); ?></span>
                  <!-- <i class="fas fa-eye pointer ml-2 text-black" onclick="togglePanelPasswd()"></i> -->
                  <!-- <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPanelPasswd" data-clipboard-action="copy" data-clipboard-target="#panelPasswd"></i> -->
                </p>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; ?>
            <!-- 面板管理信息end -->
          </div>
          <!-- 流量 -->
          <?php if(($Detail['host_data']['domainstatus'] == 'Active' || ($Detail['host_data']['domainstatus'] == 'Suspended' && $Detail['host_data']['suspendreason_type'] == 'flow')) && $Detail['host_data']['bwlimit'] > 0): ?>
            <!-- <div class="d-flex justify-content-end mb-2">
              <button type="button" class="btn btn-success btn-sm waves-effect waves-light" id="orderFlowBtn"
                onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
            </div> -->
          <div class="mt-4 mb-3">
            <i class="bx bx-circle" style="color:#f0ad4e"></i> <?php echo $Lang['used_flow']; ?>：<span id="usedFlowSpan">-</span>
            <i class="bx bx-circle" style="color:#34c38f"></i> <?php echo $Lang['residual_flow']; ?>：<span id="remainingFlow">-</span>
          </div>
          <div class="mb-4 ll-flex">
            <div class="progress w-75">
              <div class="progress-bar progress-bar-striped bg-success" id="totalProgress" role="progressbar"
                style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">100%
              </div>
            </div>
            <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb ml-2" id="orderFlowBtn"
              onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
          </div>
          <?php endif; ?>
          <!-- 流量end -->
          <div class="row">
            <!-- 订购价格start -->

            <!-- <div class="col-12 my-2">
                <div class="d-flex justify-content-between align-items-center">
                  <span><?php echo $Lang['first_order_price']; ?></span>
                  <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): if($Detail['host_data']['status'] == 'Paid'): ?>
                  <button type="button" class="btn btn-primary btn-sm waves-effect waves-light" id="renew"
                    onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['immediate_renewal']; ?></button>
                  <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                  <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                    <button type="button" class="btn btn-primary btn-sm waves-effect waves-light"
                      id="renewpay"><?php echo $Lang['immediate_renewal']; ?></button>
                  </a>
                  <?php endif; ?>
                  <?php endif; ?>
                </div>
              </div> -->
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray pre-money-box">
                <p class="text-gray">
                  <?php echo $Lang['first_order_price']; ?>
                  <span class="fr"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): else: ?>
                    <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span>
                </p>

                <section class="d-flex align-items-center">
                  <h3 class="mb-0 mr-2 dc">
                    <?php echo !empty($Detail['host_data']['firstpaymentamount_desc']) ? $Detail['host_data']['firstpaymentamount_desc'] : '-'; ?></h3>
                  <!-- <span class="badge
                        <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): ?> - <?php else: ?> <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span> -->
                </section>
                <?php if($Renew['host']['billingcycle'] != 'onetime' && $Renew['host']['status'] == 'Paid' && $Renew['host']['billingcycle'] !=
                'free'): ?>
                <section class="d-flex align-items-center flex-wrap">
                  <span><?php echo $Lang['automatic_balance_renewal']; ?></span>
                  <div class="custom-control custom-switch custom-switch-md mb-4 ml-2" dir="ltr">
                    <input type="checkbox" class="custom-control-input" id="automaticRenewal" 
                      onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew'] !=0): ?>checked
                      <?php endif; if($Detail['host_data']['billingcycle'] ==ontrial): ?>disabled
                      <?php endif; ?> > <label class="custom-control-label" for="automaticRenewal"></label>
                  </div>
                  <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): if($Detail['host_data']['status'] == 'Paid'): ?>
                  <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb" id="renew"
                    onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                  <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                  <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                    <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb"
                      id="renewpay"><?php echo $Lang['renew']; ?></button>
                  </a>
                  <?php endif; ?>
                  <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): if($Cancel['host_cancel']): ?>
                  <!-- <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn"
                        onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
                        ==
                        'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button> -->
                  <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn" data-container="body"
                    data-toggle="popover" data-placement="top" data-trigger="hover"
                    data-content="将于{<?php echo date('Y-m-d',!is_numeric($Detail['host_data']['deletedate'])? strtotime($Detail['host_data']['deletedate']) : $Detail['host_data']['deletedate']); ?>}自动删除"
                    onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['cancel_out']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-danger btn-sm rsb ml-2" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['apply_out']; ?></button>
                  <?php endif; ?>
                  <?php endif; ?>
                </section>
                <?php endif; ?>
                <section class="text-gray">
                  <p><?php echo $Lang['subscription_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></p>
                  <p><?php echo $Lang['payment_cycle']; ?>：<?php echo $Detail['host_data']['billingcycle_desc']; ?></p>
                  <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                  <p><?php echo $Lang['due_date']; ?>：<?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): ?> - <?php else: ?>
                    <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></p>
                  <?php else: ?>
                  <p><?php echo $Lang['due_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></p>
                  <?php endif; ?>
                </section>
              </div>
            </div>
            <!-- 订购价格end -->
          </div>
          <!-- 远程地址 建站解析 -->
          <div id="cloud_nat">
            <?php if($Detail['dcimcloud'] && ($Detail['dcimcloud']['nat_acl'] || $Detail['dcimcloud']['nat_web'])): ?>
            <div class="bg-primary rounded-sm d-flex flex-column justify-content-center text-white py-2 px-1 mb-2">
              <?php if($Detail['dcimcloud']['nat_acl']): ?>
              <div class="d-flex justify-content-between align-items-center">
                <span>
                  <label><?php echo $Lang['remote_address']; ?>：</label>
                  <span id="nat_aclBox"><?php echo $Detail['dcimcloud']['nat_acl']; ?></span>
                </span>
                <span>
                  <i class="bx bx-copy pointer text-white btn-copy" id="btnCopyaclBox" data-clipboard-action="copy"
                    data-clipboard-target="#nat_aclBox"></i>
                </span>
              </div>
              <?php endif; if($Detail['dcimcloud']['nat_web']): ?>
              <div class="d-flex justify-content-between align-items-center">
                <span>
                  <label><?php echo $Lang['ip_address']; ?>建站解析：</label>
                  <span id="nat_webBox"><?php echo $Detail['dcimcloud']['nat_web']; ?></span>
                </span>
                <span>
                  <i class="bx bx-copy pointer text-white btn-copy" id="btnCopywebBox" data-clipboard-action="copy"
                    data-clipboard-target="#nat_webBox"></i>
                </span>
              </div>
              <?php endif; ?>
            </div>
            <?php endif; ?>
          </div>
          <!-- 配置项 -->
          <div class="row">

            <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '5'||$item['option_type'] == '6'||$item['option_type'] == '8'): else: ?>

            <div class="col-md-12 mb-2 configuration configuration_list">

              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['name']; ?>：<?php echo $item['sub_name']; ?>" class="bg-light card-body bg-gray mg-0 row">

                <p class="text-gray col-md-6 plr-0 text-nowrap">
                  <?php echo $item['name']; ?>
                </p>

                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php if($item['option_type']===12): ?>
                  <img src="/upload/common/country/<?php echo $item['code']; ?>.png" width="20px">
                  <?php endif; ?>
                  <?php echo $item['sub_name']; ?>
                </p>

              </div>

            </div>

            <?php endif; ?>

            <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['fieldname']; ?>：<?php echo $item['value']; ?>" class="bg-light card-body bg-gray mg-0 row">
                <p class="text-gray col-md-6 plr-0 text-nowrap"><?php echo $item['fieldname']; ?></p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php echo $item['value']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>

            <?php endforeach; ?>
            <div onclick="isShowConfiguration()" class="configuration-btn-down isClick"><?php echo $Lang['view_more_information']; ?></div>
            <div class="col-12 mb-2">

              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $Lang['remarks_infors']; ?>：<?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>"
                class="bg-light card-body bg-gray mg-0 row">

                <p class="text-gray col-md-3 plr-0 text-nowrap"><?php echo $Lang['remarks_infors']; ?></p>

                <p class="mb-0 col-md-9 plr-0 text-nowrap"><?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>

              </div>

            </div>





          </div>

        </div>

      </div>

    </div>

    <div class="col-md-9">

      <div class="card">

        <div class="card-body">

          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">

            <?php if($Detail['module_chart']): ?>

            <li class="nav-item" id="chartLi">

              <a class="nav-link active" data-toggle="tab" href="#home1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-home"></i></span> -->

                <span><?php echo $Lang['charts']; ?></span>

              </a>

            </li>

            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>

            <li class="nav-item">

              <a class="nav-link <?php if(!$Detail['module_chart']): ?>active<?php endif; ?>" data-toggle="tab" href="#profile1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="far fa-user"></i></span> -->

                <span><?php echo $Lang['upgrade_downgrade']; ?></span>

              </a>

            </li>

            <?php endif; if($Detail['download_data']): ?>

            <li class="nav-item">

              <a class="nav-link <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#messages1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->

                <span><?php echo $Lang['file_download']; ?></span>

              </a>

            </li>

            <?php endif; if($Detail['host_data']['show_traffic_usage']): ?>
            <li class="nav-item" id="usedLi">

              <a class="nav-link <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#dosage" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->

                <span><?php echo $Lang['consumption']; ?></span>

              </a>

            </li>
            <?php endif; foreach($Detail['module_client_area'] as $item): ?>

            <li class="nav-item">

              <a class="nav-link" data-toggle="tab" href="#module_client_area_<?php echo $item['key']; ?>" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

                <span><?php echo $item['name']; ?></span>

              </a>

            </li>

            <?php endforeach; ?>

            <li class="nav-item">

              <a class="nav-link <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#settings1" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

                <span><?php echo $Lang['journal']; ?></span>

              </a>

            </li>

            <li class="nav-item">

              <a class="nav-link" data-toggle="tab" href="#finance" role="tab">

                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

                <span><?php echo $Lang['finance']; ?></span>

              </a>

            </li>

          </ul>



          <!-- Tab panes -->

          <div class="tab-content p-3 text-muted" style="min-height: 550px;">

            <?php if($Detail['module_chart']): ?>

            <div class="tab-pane active" id="home1" role="tabpanel">

              <style>
    .chartbox {
      width: 50%!important;
    }
    @media screen and (max-width: 1367px) {
      .chartbox {
        width: 100%!important;
      }
    }
</style>
<div class="row">
<?php foreach($Detail['module_chart'] as $key=>$item): ?>

  <div class="chartbox">
    <div class="module_chart module_chart_<?php echo $item['type']; ?>" data-type="<?php echo $item['type']; ?>" data-title="<?php echo $item['title']; ?>">
      <div style="height: 60px;">
        <h4 style="float: left;"><?php echo $item['title']; ?></h4>
        <?php if($item['select']): ?>
        <div style="width:200px;float: right;margin-right: 20px;">
          <select class="module_chart_select_<?php echo $item['type']; ?> form-control selectpicker_refresh" onchange="getChartDataFn('','<?php echo $key; ?>')">
            <?php foreach($item['select'] as $item2): ?>
            <option value="<?php echo $item2['value']; ?>"><?php echo $item2['name']; ?></option>
            <?php endforeach; ?>
          </select>
        </div>
        <?php endif; ?>
      </div>
      <div class="module_chart_date">
        <input class="form-control startTime" type="datetime-local" onchange="getChartDataFn('','<?php echo $key; ?>')">
        <span class="ml-1 mr-1"><?php echo $Lang['reach']; ?></span>
        <input class="form-control endTime mr-3" type="datetime-local" onchange="getChartDataFn('','<?php echo $key; ?>')">
      </div>

    </div>

    <div class="w-100 h-100 ">

      <div style="height: 500px" class="chart_content_box w-100" id="module_chart_<?php echo $item['type']; ?>"></div>

    </div>
  </div>
<?php endforeach; ?>

</div>

  




<script>
  // 图表tabs
  $(document).ready(function () {
    var arr = JSON.parse('<?php echo json_encode($Detail['module_chart']); ?>')
    setTimeout(function(){
      arr.forEach(function(item){
        getChartDataFn(item)
      })
    }, 0);
    
  });

  let switch_id = []
  let chartsData = []
  let timeArray = []
  let name = []
  let typeArray = []
  let myChart = null

  $('#chartLi').on('click', function () {
    setTimeout(function(){
      myChart.resize()
    }, 0);
  });


  // line
  function lineChartOption (type, xAxisData, seriesData0, seriesData1, unit, label) {
    // 硬盘IO
    const myChart = echarts.init(document.getElementById('module_chart_'+type))
    myChart.setOption({
      backgroundColor: '#fff',
      title: {
        subtext: (!xAxisData.length) ? '暂无数据' : '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 250
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        backgroundColor: '#fff',
        textStyle: {
          color: '#333',
          fontSize: 12
        },
        padding: [10, 10],
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);',
        formatter: function (params, ticket, callback) {
          // console.log('line:', params)
          const res = `<div>
                        <div>${params[0].marker} ${params[0].seriesName}：${params[0].value}${unit}</div>
                        <div>${params[1] ? params[1].marker : ''} ${params[1] ? params[1].seriesName : ''}${params[1] ? '：' : ''}${params[1] ? params[1].value : ''}${params[1] ? unit : ''}</div>
                        <div style="color: #999999;">${params[0].axisValue}</div>
                      </div>`
          return res
        }
      },
      grid: {
        left: '80',
        top: 30,
        x: 50,
        x2: 50,
        y2: 80
      },
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: [{
        offset: 15,
        type: 'category',
        boundaryGap: false,
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            type: 'dashed',
            color: '#ddd',
            width: 1
          }
        },
        // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
        //   return str.replace(' ', '\n')
        // }),
        data: xAxisData,
        // 轴刻度
        axisTick: {
          show: false
        },
        // 轴网格
        splitLine: {
          show: false
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999999'
          }
        }
      }],
      yAxis: [{
        type: 'value',
        axisTick: {
          show: false
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          formatter: '{value}' + unit,
          textStyle: {
            color: '#556677'
          }
        },
        splitLine: {
          lineStyle: {
            type: 'dashed'
          }
        }
      }],
      series: [{
        name: label[0],
        type: 'line',
        // data: [5, 12, 11, 14, 25, 16, 10, 18, 6],
        data: seriesData0,
        symbolSize: 1,
        symbol: 'circle',
        smooth: true,
        showSymbol: false,
        lineStyle: {
          width: 2
        },
        itemStyle: {
          normal: {
            color: '#75db16',
            borderColor: '#75db16'
          }
        }
      }, {
        name: label[1],
        type: 'line',
        // data: [10, 10, 30, 12, 15, 3, 7, 20, 15000],
        data: seriesData1,
        symbolSize: 1,
        symbol: 'circle',
        smooth: true,
        showSymbol: false,
        lineStyle: {
          width: 3,
          shadowColor: 'rgba(92, 102, 255, 0.3)',
          shadowBlur: 10,
          shadowOffsetY: 20
        },
        itemStyle: {
          normal: {
            color: '#5c66ff',
            borderColor: '#5c66ff'
          }
        }
      }
      ]
    })

    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }
  // area
  function areaChartOption (type, xAxisData, seriesData, unit, label) {
    // CPU使用率
    const myChart = echarts.init(document.getElementById('module_chart_'+type ))
    myChart.setOption({
      grid: {
        left: '80',
        top: 30,
        x: 50,
        x2: 50,
        y2: 80
      },
      backgroundColor: '#fff',
      title: {
        subtext: (!xAxisData.length) ? '暂无数据' : '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 250
        }
      },
      tooltip: {
        backgroundColor: '#fff',
        padding: [10, 20, 10, 8],
        textStyle: {
          color: '#333',
          fontSize: 12
        },
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        formatter: function (params, ticket, callback) {
          // console.log(params, '')
          const res = `<div>
                        <div>${params[0].seriesName}：${params[0].value}${unit} </div>
                        <div style="color: #999999;">${params[0].axisValue}</div>
                      </div>`
          return res
        },
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
      },
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: {
        offset: 15,
        type: 'category',
        boundaryGap: false,
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            color: '#999999',
            width: 1
          }
        },
        // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
        //   return str.replace(' ', '\n')
        // }),
        data: xAxisData,
        // 轴刻度
        axisTick: {
          show: false
        },
        // 轴网格
        splitLine: {
          show: false
        },
        axisLabel: {
          show: true,
          // interval: 0, // 横轴信息全部显示
          textStyle: {
            color: '#999999'
          }
        }
      },
      yAxis: {
        axisTick: {
          show: false // 轴刻度不显示
        },
        max: 100,
        min: 0,
        // 改变y轴颜色
        axisLine: {
          show: false
        },
        // 轴网格
        splitLine: {
          show: true,
          lineStyle: {
            color: '#ddd',
            type: 'dashed'
          }
        },
        // 坐标轴文字样式
        axisLabel: {
          show: true,
          formatter: '{value}' + unit,
          textStyle: {
            color: '#999999'
          }
        }
      },
      series: [{
        name: label,
        type: 'line',
        areaStyle: {
          opacity: 1,
          color: '#737dff'
        },
        symbol: 'none', // 折线无拐点
        lineStyle: {
          normal: {
            width: 0 // 折线宽度
          }
        },
        smooth: true,
        // data: [5, 25, 20, 50, 10, 40, 18, 25, 0]
        data: seriesData

      }]
    })

    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }

  // bar
  function barChartOption (type, xAxisData, seriesData0, seriesData1, unit, label) {
    // 内存用量
    const myChart = echarts.init(document.getElementById('module_chart_'+type))
    myChart.setOption({
      backgroundColor: '#fff',
      title: {
        subtext: (!xAxisData.length) ? '暂无数据' : '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 250
        }
      },
      tooltip: {
        backgroundColor: '#fff',
        padding: [10, 20, 10, 8],
        textStyle: {
          color: '#000',
          fontSize: 12
        },
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        formatter: function (params, ticket, callback) {
          // console.log('bar:', params)
          const res = `
          <div>
              <div>${params[0].marker}${params[0].seriesName}：${params[0].value}${unit} </div>                
              <div>${params[1] ? params[1].marker : ''} ${params[1] ? params[1].seriesName : ''}${params[1] ? '：' : ''}${params[1] ? params[1].value : ''}${params[1] ? unit : ''}</div>
              <div>${params[0].axisValue}</div>
          </div>`
          return res
        },
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
      },
      grid: {
        left: '80',
        top: 30,
        x: 70,
        x2: 50,
        y2: 80
      },
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: {
        offset: 15,
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          }
        },
        type: 'category',
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            type: 'dashed',
            color: '#ddd',
            width: 1
          }
        },
        // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
        //   return str.replace(' ', '\n')
        // })
        data: xAxisData
      },
      yAxis: {
        axisTick: {
          show: false // 轴刻度不显示
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          },
          formatter: '{value}' + unit
        },
        // 轴网格
        splitLine: {
          show: true,
          lineStyle: {
            color: '#ddd',
            type: 'dashed'
          }
        }

      },
      series: [{
        name: label[1],
        type: 'bar',
        stack: '总量',
        barGap: '-100%',
        // data: [136, 132, 101, 134, 90, 230, 210, 100, 300],
        data: seriesData1,
        itemStyle: {
          barBorderRadius: [5, 5, 0, 0],
          color: '#737dff'
        }
      },
      {
        name: label[0],
        type: 'bar',
        stack: '可用',
        // data: [964, 182, 191, 234, 290, 330, 310, 100, 500],
        data: seriesData0,
        itemStyle: {
          barBorderRadius: [5, 5, 0, 0],
          color: '#ccc',
          opacity: 0.3
        }
      } 
      ]
    })

    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }



  async function getChartDataFn(e,index) {
    

    var arr = JSON.parse('<?php echo json_encode($Detail['module_chart']); ?>')
    if(index){e = arr[index]}
    var echartTT = echarts.init(document.getElementById('module_chart_'+e.type))
    echartTT.showLoading({
      text: '数据正在加载...',
      color: '#999',
      textStyle: {
        fontSize: 30,
        color: '#444'
      },
      effectOption: {
        backgroundColor: 'rgba(0, 0, 0, 0)'
      }
    })
    const queryObj = {
      id: '<?php echo app('request')->get('id'); ?>',
      type: e.type,
      start: new Date($('.module_chart_'+e.type+' .startTime').val()).getTime(),
      end:new Date($('.module_chart_'+e.type+' .endTime').val()).getTime(),
      select: $('.module_chart_select_' + e.type).val(),
    }
    $.ajax({
      type: "GET",
      url: '' + '/provision/chart/<?php echo app('request')->get('id'); ?>',
      data: queryObj,
      success: function (data) {
        echartTT.hideLoading()
      if (data.status !== 200) return false

      const xAxisData = []
      const seriesData0 = []
      const seriesData1 = [];

      (data.data.list || []).forEach((item, index) => {
        (item || []).forEach(innerItem => {
          if (index === 0) {
            xAxisData.push(innerItem.time)
            seriesData0.push(innerItem.value)
          } else if (index === 1) {
            seriesData1.push(innerItem.value)
          }
        })
      })

       if (data.data.chart_type === 'area') {
          areaChartOption(e.type, xAxisData, seriesData0, data.data.unit, data.data.label)
        } else if (data.data.chart_type === 'line') {
          lineChartOption(e.type, xAxisData, seriesData0, seriesData1, data.data.unit, data.data.label)
        } else if (data.data.chart_type === 'bar') {
          barChartOption(e.type, xAxisData, seriesData0, seriesData1, data.data.unit, data.data.label)
        }

        // 如果初始查询没有时间, 则设置默认时间为返回数据的第一个和最后一个时间
        if (!$('.module_chart_'+e.type+' .startTime').val() || !$('.module_chart_'+e.type+' .endTime').val()) {
          if (data.data.list[0].length) {
            var start = new Date(data.data.list[0][0].time).getTime()
            var end = new Date(data.data.list[0][data.data.list[0].length - 1].time).getTime()
            $('.module_chart_'+e.type+' .startTime').val(moment(start).format('YYYY-MM-DDTHH:mm'))
            $('.module_chart_'+e.type+' .endTime').val(moment(end).format('YYYY-MM-DDTHH:mm'))
          }
        }
      }
    });
  }





</script>

            </div>

            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <div class="tab-pane <?php if(!$Detail['module_chart']): ?>active<?php endif; ?>" id="profile1" role="tabpanel">

              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">

                  <div class="col-12">

                    <div class="bg-light  rounded card-body">

                      <div class="row">

                        <div class="col-sm-3">
                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>

                        <div class="col-sm-6">



                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>

                        </div>

                        <div class="col-sm-3">

                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>

                        </div>

                      </div>

                    </div>

                  </div>

                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">

                  <div class="col-12">

                    <div class="bg-light  rounded card-body">

                      <div class="row">

                        <div class="col-sm-3">

                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>

                        </div>

                        <div class="col-sm-6">

                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>

                        </div>

                        <div class="col-sm-3">

                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>

                        </div>

                      </div>

                    </div>

                  </div>

                </div>
                <?php endif; ?>

              </div>

            </div>
            <?php endif; if($Detail['download_data']): ?>

            <div
              class="tab-pane <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
              id="messages1" role="tabpanel">

              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>

            </div>

            <?php endif; ?>

            <div
              class="tab-pane <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
              id="settings1" role="tabpanel">

              <!-- 日志 -->

            </div>

            <?php foreach($Detail['module_client_area'] as $item): ?>

            <div class="tab-pane" role="tabpanel" id="module_client_area_<?php echo $item['key']; ?>">
              <div style="width:100%;min-height: 550px;">
                <script>
                    $.ajax({
                       url : '/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>'
                      ,type : 'get'
                      ,success : function(res)
                      {
                          $('#module_client_area_<?php echo $item['key']; ?> > div').html(res);
                      }
                    })
                </script>
              </div>

              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                onload="this.height=$($('.main-content .card-body')[1]).height()-72" frameborder="0"
                width="100%"></iframe> -->
              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                frameborder="0" width="100%" style="min-height: 550px;"></iframe> -->

            </div>

            <?php endforeach; if($Detail['host_data']['show_traffic_usage']): ?>
            <div
              class="tab-pane <?php if(!$Detail['module_chart'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product'] && !$Detail['download_data']): ?>active<?php endif; ?>"
              id="dosage" role="tabpanel">

              <div class="row d-flex align-items-center">

                <div class="col-md-3">

                  <input class="form-control" type="date" id="startingTime">

                </div>

                <span><?php echo $Lang['reach']; ?></span>

                <div class="col-md-3">

                  <input class="form-control" type="date" id="endTime">

                </div>

              </div>

              <div class="w-100 h-100">

                <div style="height: 500px" class="chart_content_box w-100" id="usedChartBox"></div>

              </div>

            </div>
            <?php endif; ?>

            <div class="tab-pane" id="finance" role="tabpanel">

              <!-- 财务 -->

            </div>

          </div>

        </div>

      </div>

    </div>

  </div>

</div>

<!-- 魔方云救援系统弹窗 --> 
<div class="modal fade" id="moduleDcimCloudRescue" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
     aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title mt-0">救援系统</h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <form id="dcimcloudRescue">
                <div class="modal-body">

                    <div class="form-group row mb-4">
                        <label for="horizontal-firstname-input" class="col-md-3 col-form-label d-flex justify-content-end"><span style="color: red;"> * </span>救援系统类型</label>
                        <div class="col-md-6">
                            <select class="form-control" name="system">
                                <option value="1">Windows</option>
                                <option value="2">Linux</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group row mb-4">
                        <label for="horizontal-firstname-input" class="col-md-3 col-form-label d-flex justify-content-end"><span style="color: red;"> * </span>临时密码</label>
                        <div class="col-md-6">
                            <input type="text" class="form-control" name="temp_pass">
                        </div>
                        <div class="col-md-1 fs-18 d-flex align-items-center">
                            <i class="fas fa-dice create_random_pass pointer" onclick="create_random_pass()"></i>
                        </div>
                        <label id="password-error-tip" class="control-label error-tip" for="temp_pass"></label>
                    </div>
                    <div class="alert alert-danger" role="alert">
                        请妥善保存当前密码，该密码不会二次使用
                    </div>
                    
                    <div id="moduleDcimCloudRescueForceDiv" style="display: none;">
                      <div class="alert alert-danger" role="alert">
                      当前操作需要实例在关机状态下进行<br>
                      为了避免数据丢失，实例将关机中断您的业务，请仔细确认。<br>
                      强制关机可能会导致数据丢失或文件系统损坏，您也可以主动关机后再进行操作。
                      </div>
                      
                      <div class="form-group row mb-4">
                          <label for="horizontal-firstname-input" class="col-md-3 col-form-label d-flex justify-content-end">强制关机</label>
                          <div class="col-md-6 pt-9">
                              <label>
                                  <input type="checkbox" class="mr-1" name="force">同意强制关机
                              </label>
                          </div>
                          <label class="control-label error-tip force-error-tip" for="force">请同意强制关机</label>
                      </div>
                    </div>
                </div>
                <input type="hidden" name="func" value="rescue_system">
                <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
            </form>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
                <button type="button" class="btn btn-primary submit" onclick="moduleDcimCloudRescue($(this))"><?php echo $Lang['determine']; ?></button>
            </div>
        </div>
    </div>
</div>

<!-- 电源状态 -->
<script>
  var showPowerStatus = '<?php echo $Detail['module_power_status']; ?>';
  var powerStatus = {
    status: '',
    des: ''
  }

  $(function () {
    $(".nav-tabs-custom").find('.nav-item').find("a").removeClass('active')
    $(".nav-tabs-custom").find('.nav-item').eq(0).find("a").addClass('active')
    // 暂停状态悬浮原因
    $('.container-fluid').on('mouseover', '.status-suspended', function () {
      $('.xf-bg').show();
    })
    $('.container-fluid').on('mouseout', '.status-suspended', function () {
      $('.xf-bg').hide();
    })

    if (showPowerStatus == '1') {
      getPowerStatus('<?php echo app('request')->get('id'); ?>')
    }

    $('#powerBox').on('click', function () {
      getPowerStatus('<?php echo app('request')->get('id'); ?>')
    });
  })

  // 获取电源状态
  function getPowerStatus() {

    $('#powerStatusIcon').removeClass()
    $('#powerStatusIcon').addClass('bx bx-loader')

    $.ajax({
      type: "POST",
      url: '' + '/provision/default',
      data: {
        id: '<?php echo app('request')->get('id'); ?>',
        func: 'status'
      },
      success: function (data) {
        $('#powerStatusIcon').attr('data-content', data.data ? data.data.des : data.msg)
        $('#powerStatusIcon').removeClass()
        if (data.status != 200) {
          powerStatus.status = 'unknown'
          powerStatus.des = data.msg
          $('#powerStatusIcon').addClass('sprite unknown')
        } else {
          powerStatus.status = data.data ? data.data.status : 'unknown'
          powerStatus.des = data.data ? data.data.des : '<?php echo $Lang['unknown']; ?>'
           
          $('#moduleDcimCloudRescueForceDiv').hide();
          $('#moduleDcimCloudRescueForceDiv').removeClass('force_show');
          if (powerStatus.status === 'process') {
            $('#powerStatusIcon').addClass('bx bx-loader')
            $('#moduleDcimCloudRescueForceDiv').show();
            $('#moduleDcimCloudRescueForceDiv').addClass('force_show');
          } else if (powerStatus.status === 'on') {
            $('#powerStatusIcon').addClass('sprite start')
            $('#moduleDcimCloudRescueForceDiv').show();
            $('#moduleDcimCloudRescueForceDiv').addClass('force_show');
          } else if (powerStatus.status === 'off') {
            $('#powerStatusIcon').addClass('sprite closed')
          } else if (powerStatus.status === 'waiting') {
            $('#powerStatusIcon').addClass('sprite waitOn')
          } else if (powerStatus.status === 'suspend') {
            $('#powerStatusIcon').addClass('sprite pause')
          } else if (powerStatus.status === 'wait_reboot' || powerStatus.status === 'wait') {
            $('#powerStatusIcon').addClass('sprite waiting')
          } else if (powerStatus.status === 'cold_migrate') {
            $('#powerStatusIcon').addClass('iconfont icon-shujuqianyi')
          } else if (powerStatus.status === 'hot_migrate') {
            $('#powerStatusIcon').addClass('iconfont icon-shujuqianyi')
          } else {
            $('#powerStatusIcon').addClass('sprite unknown')
          }

          if (data.data.status !== 'process') { // 状态改变, 清除定时器
            clearInterval(timeInterval)
          }
        }

      }
    });
  }

  var timeOut = null
  var timeInterval = null
  //点击更多信息
  function isShowConfiguration(first = true) {
    let time = 0
    if (first) {
      time = 500
    }
    for (let i = 0; i < 99; i++) {
      if (i > 3) {
        $('.configuration').eq(i).toggle(time)
      }
    }
    $('.configuration-btn-down').toggleClass('isClick')
    $('.configuration-btn-down').html('查看更多信息')
    $('.configuration-btn-down.isClick').html('收起更多信息')
    if($('.configuration_list').length < 4) {
      $('.configuration-btn-down').hide()
    }
  }
  isShowConfiguration(false)
</script>
<script>
  // var showPWd = true
  $(function () {
    // 复制IP
    var clipboard1 = new ClipboardJS('#btnCopyIP', {
      text: function (trigger) {
        return $('#copyIPContent').text()
        // return $("#copyIPOne").text()
      }
    });
    clipboard1.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })

    // 一个ip时，不弹框，复制ip
    var clipboard = new ClipboardJS('.copyOneIp', {
      text: function (trigger) {
        return $('#copyOneIp').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })

    // 查看密码
    // $('#copyPwdContent').hide()
    // $('#hidePwdBox').hide()

    // 复制密码
    var clipboard2 = new ClipboardJS('.btnCopyPwd', {
      text: function (trigger) {
        return $('#copyPwdContent').text()
      }
    });
    clipboard2.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })
  // function togglePwd() {
  //   showPWd = !showPWd

  //   if (showPWd) {
  //     $('#copyPwdContent').show()
  //     $('#hidePwdBox').hide()
  //   }
  //   if (!showPWd) {
  //     $('#copyPwdContent').hide()
  //     $('#hidePwdBox').show()
  //   }
  // }

  function togglePanelPasswd() {
    if ($("#hidePanelPasswd").is(':hidden')) {
      $("#hidePanelPasswd").show();
      $("#panelPasswd").hide();
    } else {
      $("#hidePanelPasswd").hide();
      $("#panelPasswd").show();
    }
  }

  // 远程地址 建站解析
  // 复制授权码
  var clipboardnat_aclBox = new ClipboardJS('#btnCopyaclBox', {
    text: function (trigger) {
      return $('#nat_aclBox').text()
    }
  });
  clipboardnat_aclBox.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

  // 复制授权码
  var clipboardnat_webBox = new ClipboardJS('#btnCopywebBox', {
    text: function (trigger) {
      return $('#nat_webBox').text()
    }
  });
  clipboardnat_webBox.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

  // 复制面板管理密码
  var clipboard_panelPasswd = new ClipboardJS('.btnCopyPanelPasswd', {
    text: function (trigger) {
      return $('#panelPasswd').text()
    }
  });
  clipboard_panelPasswd.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })
</script>


<script>
  // 用量tabs

  let usedChart = null

  let usedStartTime

  let usedEndTime



  $(document).ready(function () {
    //chartOption()

    if ($('#startingTime,#endTime').length > 0) getData()

    window.addEventListener('resize', function () {

      if (usedChart) usedChart.resize()

    })



    $('#usedLi').on('click', function () {

      if (usedChart) {
        setTimeout(function () {

          usedChart.resize()

        }, 0);
      }

    });


    $('#startingTime,#endTime').change(function () {

      usedStartTime = $('#startingTime').val()

      usedEndTime = $('#endTime').val()

      getData()

    });

    // 获取数据

    async function getData() {
      usedChart = echarts.init(document.getElementById('usedChartBox'))

      usedChart.showLoading({

        text: '数据正在加载...',

        color: '#999',

        textStyle: {

          fontSize: 30,

          color: '#444'

        },

        effectOption: {

          backgroundColor: 'rgba(0, 0, 0, 0)'

        }

      })

      const obj = {

        id: '<?php echo app('request')->get('id'); ?>',

        start: usedStartTime,

        end: usedEndTime

      }

      $.ajax({

        type: "get",

        url: '/host/trafficusage',

        data: obj,

        success: function (data) {

          usedChart.hideLoading()

          if (data.status !== 200) return false

          const xAxisData = []

          const seriesData = [];
          const seriesData2 = [];

          (data.data || []).forEach(item => {

            xAxisData.push(item.time)
            seriesData.push(item.in)
            seriesData2.push(item.out)

          })

          usedChart = echarts.init(document.getElementById('usedChartBox'))

          usedChart.setOption({
            backgroundColor: '#fff',
            title: {
              subtext: '',
              left: 'center',
              textAlign: 'left',
              subtextStyle: {
                lineHeight: 250
              }
            },
            tooltip: {
              backgroundColor: '#fff',
              padding: [10, 20, 10, 8],
              textStyle: {
                color: '#000',
                fontSize: 12
              },
              trigger: 'axis',
              axisPointer: {
                type: 'line',
                lineStyle: {
                  color: '#7dcb8f'
                }
              },
              formatter: function (params, ticket, callback) {
                //console.log('bar:', params)
                const res = `
          <div>
              <div>接收流量(GB)：${params[0].value} </div>                
              <div>发送流量(GB)：${params[1].value}</div>
              <div>${params[0].axisValue}</div>
          </div>`
                return res
              },
              extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
            },
            grid: {
              left: '80',
              top: 30,
              x: 70,
              x2: 50,
              y2: 80
            },
            dataZoom: [ // 缩放
              {
                type: 'inside',
                throttle: 50
              }
            ],
            xAxis: {
              offset: 15,
              axisLabel: {
                show: true,
                textStyle: {
                  color: '#999'
                }
              },
              type: 'category',
              // 改变x轴颜色
              axisLine: {
                lineStyle: {
                  type: 'dashed',
                  color: '#ddd',
                  width: 1
                }
              },
              // data: ['2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00', '2020-08-11 11:30:00'].map(function (str) {
              //   return str.replace(' ', '\n')
              // })
              data: xAxisData
            },
            yAxis: {
              axisTick: {
                show: false // 轴刻度不显示
              },
              axisLine: {
                show: false
              },
              axisLabel: {
                show: true,
                textStyle: {
                  color: '#999'
                },
                //  formatter: '{value}' + 'GB'
              },
              // 轴网格
              splitLine: {
                show: true,
                lineStyle: {
                  color: '#ddd',
                  type: 'dashed'
                }
              }

            },
            series: [{
                name: '21',
                type: 'bar',
                stack: '接收流量',
                barGap: '-100%',
                //data: [136, 132, 101, 134, 90, 230, 210, 100, 300],
                data: seriesData,
                itemStyle: {
                  barBorderRadius: [5, 5, 0, 0],
                  color: '#737dff'
                }
              },
              {
                name: '2',
                type: 'bar',
                stack: '使用流量',
                // data: [964, 182, 191, 234, 290, 330, 310, 100, 500],
                data: seriesData2,
                itemStyle: {
                  barBorderRadius: [5, 5, 0, 0],
                  color: '#ccc',
                  opacity: 0.3
                }
              }
            ]
          })

          // 如果初始查询没有时间, 则设置默认时间为返回数据的第一个和最后一个时间

          if (!usedStartTime || !usedEndTime) {

            if (data.data.length) {

              usedStartTime = new Date().getFullYear() + '-' + data.data[0].time

              usedEndTime = new Date().getFullYear() + '-' + data.data[data.data.length - 1].time

              $('#startingTime').val(data.data[0].time);

              $('#endTime').val(data.data[data.data.length - 1].time);

            }

          }

        }

      });

    }



    // 分辨率改变, 重绘图表

    function resize() {

      usedChart.resize()

    }



    // 时间选择改变

    function dateChange() {

      const startTimeStamp = new Date(usedStartTime).getTime()

      const endTimeStamp = new Date(usedEndTime).getTime()

      if (usedStartTime && usedEndTime && startTimeStamp < endTimeStamp) {

        getData()

      }

    }
  })
</script>



<script>
  // 获取基础数据

  const obj = {

    host_id: '<?php echo app('request')->get('id'); ?>'

  }

  $.ajax({

    type: "get",

    url: '' + '/host/dedicatedserver',

    data: obj,

    success: function (data) {

      const totalFlow = data.data.host_data.bwlimit // 总流量

      const usedFlow = data.data.host_data.bwusage.toFixed(1) // 已用流量

      const remainingFlow = (totalFlow - usedFlow).toFixed(1)

      let percentUsed = 100 - parseInt((usedFlow / totalFlow) * 100) || 0

      $('#totalProgress')
        .css('width', percentUsed + '%')
        .attr('aria-valuenow', percentUsed)
        .text(`${percentUsed}%`);

      $('#usedFlowSpan').text(`${usedFlow > 1024 ? ((usedFlow / 1024).toFixed(2) + 'TB') : (usedFlow + 'GB')}`);
      $('#remainingFlow').text(
        `${remainingFlow > 1024 ? ((remainingFlow / 1024).toFixed(2) + 'TB') : (remainingFlow + 'GB')}`);

      // 产品状态
      $('#statusBox').append(`<span class="sprite2 ${data.data.host_data.domainstatus}"></span>`)

      if(typeof data.data.host_data.os !== 'undefined' && data.data.host_data.os.indexOf('win') == -1 && data.data.host_data.os.indexOf('Win') == -1){
        $('#moduleDcimCloudRescue select[name="system"]').find('option').eq(1).attr('selected', 'selected');
      }
    }

  });
</script>



<script>
  const logObj = {

    id: '<?php echo app('request')->get('id'); ?>',

    action: 'log_page'

  }

  $.ajax({

    type: "get",

    url: '' + '/servicedetail',

    data: logObj,

    success: function (data) {

      $(data).appendTo('#settings1');

    }

  });

  // 财务的append

  const financeObj = {

    id: '<?php echo app('request')->get('id'); ?>',

    action: 'billing_page'

  }

  $.ajax({

    type: "get",

    url: '' + '/servicedetail',

    data: financeObj,

    success: function (data) {

      $(data).appendTo('#finance');

    }

  });
</script>


<script>
  var clipboard = null
  var clipboardpoppwd = null
  var ips = <?php echo json_encode($Detail['host_data']['assignedips']); ?>;
  $(document).on('click', '#copyIPContent', function () {
    // if (ips.length<1) {
    //     return
    // }
    $('#popModal').modal('show')
    $('#popTitle').text('IP地址')
    var iplist = ''
    if (clipboard) {
      clipboard.destroy()
    }
    for(let item in ips) {
      iplist += `
        <div>
          <span class="copyIPContent${item}">${ips[item]}</span>
          <i class="bx bx-copy pointer text-primary ml-1 btn-copy btnCopyIP${item}" data-clipboard-action="copy" data-clipboard-target=".copyIPContent${item}"></i>
        </div>
      `

      // 复制IP
      clipboard = new ClipboardJS('.btnCopyIP'+item, {
        text: function (trigger) {
          return $('.copyIPContent'+item).text()
        },
        container: document.getElementById('popModal')
      });
      clipboard.on('success', function (e) {
        toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
      })
    }

    $('#popContent').html(iplist)
  });


  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      <div>
        <?php if(in_array('crack_pass', array_column($Detail['module_button']['control'], 'func')) && $Detail['host_data']['domainstatus']=="Active"): ?>
        <button type="button" class="service_module_button ml-2 mt-2 left_wrap_btn fs-12" data-func="crack_pass"
          data-type="default"
          onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"><?php echo $Lang['reset_password']; ?></button>
        <?php endif; ?>
      </div>
    `)
  });


  $('#popModal').on('shown.bs.modal', function () {
    if (clipboardpoppwd) {
      clipboardpoppwd.destroy()
    }
    clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })

  //net
  var url = setting_web_url + '/servicedetail?action=nat&id=<?php echo app('request')->get('id'); ?>';
  $.ajax({
    url: url,
    type: 'GET',
    success: function (res) {
      $('#cloud_nat').html(res)
    }
  })

  // 获取remoteinfo
  var url = setting_web_url + '/provision/custom/<?php echo app('request')->get('id'); ?>';
  $.ajax({
    url: url,
    type: 'POST',
    data: {
        func: 'remoteInfo',
    },
    success: function (res) {
        if(res.status==200){
            if(res.data.rescue == 0){
                $('.service_module_button').eq(7).hide();
            }else{
                $('.service_module_button').eq(6).hide();
            }
        }else{
            $('.service_module_button').eq(7).hide();
        }
    }
  })


</script>
<?php elseif($Detail['host_data']['type'] == "dcim"): ?>
    <style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<style>
  .server_header_box {
    height: auto;
    background-image: linear-gradient(87deg, #4d83ff 0%, #3656ff 100%);
    border-radius: 15px;
    padding: 20px 25px;
    color: #ffffff;
  }

  .left_wrap_btn {
    display: inline-block;
    width: 80px;
    height: 20px;
    background-color: #5f88fe;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    color: #ffffff;
    text-align: center;
    border: none;
  }

  .custom-button {
    background-color: #6f87fc;
    box-shadow: 0px 6px 14px 2px rgba(6, 31, 179, 0.26);
    border-radius: 4px;
    font-size: 12px;
    color: #fff;
    border: none;
  }

  .box_left_wrap {
    border-left: 1px solid rgba(255, 255, 255, 0.25);
    min-height: 74px;
  }

  .aibiao a {
    width: 100%;
    height: 100%;
    display: inline-block;
  }

  @media screen and (max-width: 1367px) {
    .form-control {
      width: 46%;
    }

    .server_header_box {
      height: auto;
    }

    .power_box {
      max-width: 300px;
    }

    .left_wrap_btn {
      width: 60px !important;
    }

    .bottom-box {
      margin-top: 3rem !important;
    }

    .osbox {
      max-width: 150px;
    }
  }

  @media screen and (max-width: 976px) {
    .server_header_box {
      height: auto;
      padding: 20px;
      margin-top: 10px;
    }

    .domain,
    .box_left_wrap {
      margin-bottom: 20px;
      border-left: none;
    }

    .power_box {
      margin-bottom: 20px;
    }
  }

  .tuxian {
    cursor: pointer;
  }

  .tuxian:hover {
    color: rgba(224, 224, 224, 0.877);
  }

  .alarm {
    display: inline-block;
    font-size: 12px;
    cursor: pointer;
    color: #495057;
    font-weight: 300;
  }

  .fr {
    float: right;
  }

  .restall-btn {
    border-radius: 25px;
    margin-left: 20px;
  }

  .login-info-icon {
    color: #5f88fe;
  }

  .dc {
    color: #5f88fe
  }

  .rsb {
    height: 20px;
    padding: 0px 10px;
  }

  .mg-0 {
    margin: 0;
  }

  .plr-0 {
    padding-left: 0px;
    padding-right: 0px;
    margin-bottom: 0;
  }

  #copyIPContent:hover {
    color: #FCA426
  }

  #copyOneIp:hover {
    color: #FCA426
  }

  .text-nowrap {
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .text-right {
    text-align: right;
  }

  .pre-money-box {
    background: url("/themes/clientarea/default/assets/images/money.png") no-repeat;
    background-position-x: right;
    background-position-y: bottom;
  }

  .w-75 {
    width: 75% !important;
  }

  .ll-flex {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .xf-bg {
    display: none;
    position: absolute;
    background: #fff;
    padding: 10px 15px;
    border-radius: 4px;
    top: -40px;
    left: 0px;
    box-shadow: 0px 3px 5px 0px rgba(0, 28, 144, 0.21);
    font-size: 12px;
  }

  .xf-bg-text {
    color: #333;
    word-break: break-all;
  }

  .flex-wrap {
    display: flex !important;
    flex-flow: wrap;
  }

  .configuration-btn-down {
    width: 100%;

    text-align: center;
    line-height: 36px;
    color: #5F88FE;
  }

  .configuration-btn-down::after {
    content: '';
    display: inline-block;
    width: 8px;
    height: 8px;
    margin-left: 8px;
    background-color: transparent;
    transform: rotate(225deg);
    border: 1px solid #5F88FE;
    border-bottom: none;
    border-right: none;
    transform-origin: 2px;
    transition: all .2s;
  }

  .configuration-btn-down.isClick::after {
    transform: rotate(45deg);
  }
</style>
<script src="/themes/clientarea/default/assets/libs/echarts/echarts.min.js?v=<?php echo $Ver; ?>"></script>
<div class="container-fluid">
  <div class="row mb-4">

    <div class="col-12">
      <div class="row align-items-center server_header_box">
        <div class="mr-3 power_box">
          <div class="text-white d-flex">
            <!-- 电源状态 -->
            <div class="mr-3 pointer">
              <?php if($Detail['module_power_status'] == '1'): ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="powerBox">
                <span id="powerStatusIcon" class="bx bx-loader" data-toggle="popover" data-trigger="hover" title=""
                  data-html="true" data-content="<?php echo $Lang['loading']; ?>..."></span>
              </div>
              <?php else: ?>
              <div class="powerimg d-flex justify-content-center align-items-center" id="statusBox"></div>
              <?php endif; ?>
            </div>
            <div>
              <section class="d-flex align-items-center mb-2">
                <h4 class="text-white mb-0 font-weight-bold"><?php echo $Detail['host_data']['productname']; ?></h4>
                <span class="badge badge-pill ml-2 py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>"
                  style="position: relative;">
                  <div class="xf-bg">
                    <div class="xf-bg-text"><span style="color: #e31519;">暂停原因：</span><?php echo !empty($Detail['host_data']['suspendreason']) ? $Detail['host_data']['suspendreason'] : $Detail['host_data']['suspendreason_type']; ?></div>
                    <font class="sj"></font>
                  </div>
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>
                </span>
              </section>
              <section>
                <span><?php echo $Detail['host_data']['domain']; ?></span>
                <span class="cancelBtn" id="cancelDcimTask" style="display:none;"
                  onclick="cancelDcimTask('<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['cancel_task']; ?></span>
              </section>
            </div>
          </div>
        </div>
        <div class="pl-4 mr-3 box_left_wrap osbox">
          <span class="text-white-50 fs-12"><?php echo $Lang['operating_system']; ?></span>
          <?php if($Detail['host_data']['domainstatus']=="Active"): ?>
          <span class="ml-0">
            <button type="button" class="dcim_service_module_button left_wrap_btn fs-12 restall-btn"
              data-func="reinstall" data-type="default"
              onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['reinstall_system']; ?></button>
          </span>
          <?php endif; ?>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $Detail['host_data']['os']; ?></h5>
        </div>

        <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '6'||$item['option_type'] == '8'): ?>

        <div class="pl-4 mr-3 box_left_wrap">
          <span class="text-white-50 fs-12"><?php echo $item['name']; ?></span>
          <h5 class="mt-2 font-weight-bold text-white"><?php echo $item['sub_name']; ?></h5>
        </div>

        <?php endif; ?>

        <?php endforeach; ?>

        <div class="pl-4 mr-3 box_left_wrap">

          <span class="text-white-50 fs-12"><?php echo $Lang['ip_address']; ?></span>

          <h5 class="mt-2 font-weight-bold text-white">
            <!-- <span data-toggle="popover" data-trigger="hover" title="" data-html="true" data-content="
              <?php foreach($Detail['host_data']['assignedips'] as $list): ?>
              <div><?php echo $list; ?></div>
              <?php endforeach; ?>
            "> -->
            <span>
              <?php if($Detail['host_data']['dedicatedip']): if($Detail['host_data']['assignedips']): ?>
              <!-- <span class="tuxian"><?php echo $Detail['host_data']['dedicatedip']; ?></span>-->
              <span id="copyIPContent" class="pointer"><?php echo $Detail['host_data']['dedicatedip']; ?>
                (<?php echo count($Detail['host_data']['assignedips']); ?>)</span>
              <?php else: ?>
              <span id="copyOneIp" class="pointer copyOneIp"><?php echo $Detail['host_data']['dedicatedip']; ?></span>
              <?php endif; ?>
              <!-- <?php if(count($Detail['host_data']['assignedips']) >= 0): ?>
            <i class="bx bx-copy pointer text-white ml-1 btn-copy" id="btnCopyIP" data-clipboard-action="copy"
              data-clipboard-target="#copyIPOne"></i>
            <?php endif; ?> -->
              <?php else: ?>
              -
              <?php endif; ?>
            </span>

            <!-- <?php if(count($Detail['host_data']['assignedips']) > 0): ?>
        <div class="alarm" id="copyIPContent">
          更多
        </div>
        <?php endif; ?> -->

          </h5>


        </div>

        <!--
      <div class="pl-4 mr-3 box_left_wrap">
      <span class="text-white-50 fs-12"><?php echo $Lang['password']; ?></span>
      <h5 class="mt-2" data-toggle="popover" data-trigger="hover" data-html="true"
        data-content="<?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?><br><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?>">
        <span id="hidePwdBox" class="text-white">***********</span>
        <span id="copyPwdContent" class="text-white"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="fas fa-eye pointer ml-2 text-white" onclick="togglePwd()"></i>
        <i class="bx bx-copy pointer ml-1 btn-copy text-white" id="btnCopyPwd" data-clipboard-action="copy"
          data-clipboard-target="#copyPwdContent"></i>
      </h5>

    </div>
    -->
        <div class="d-flex justify-content-end flex-shrink-1 flex-grow-1">


          <?php if($Detail['host_data']['type'] == "dcim" && $Detail['dcim']['auth'] && $Detail['host_data']['domainstatus']=="Active"): ?>
          <div class="btn-group ml-2 mr-2 mt-2">
            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i class="mdi mdi-chevron-down"></i></button>
            <div class="dropdown-menu">
              <?php if($Detail['dcim']['auth']['on'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="on"
                data-des="<?php echo $Lang['confirm_turn_on']; ?>"><?php echo $Lang['batch_operation']; ?></a>
              <?php endif; if($Detail['dcim']['auth']['off'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="off"
                data-des="<?php echo $Lang['confirm_turn_off']; ?>"><?php echo $Lang['shut_down']; ?></a>
              <?php endif; if($Detail['dcim']['auth']['reboot'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="reboot"
                data-des="<?php echo $Lang['confirm_turn_restart']; ?>"><?php echo $Lang['restart']; ?></a>
              <?php endif; if($Detail['dcim']['auth']['bmc'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="bmc"
                data-des="<?php echo $Lang['confirm_turn_mbc']; ?>"><?php echo $Lang['reset_bmc']; ?></a>
              <?php endif; if($Detail['dcim']['auth']['rescue'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')"
                data-func="rescue"><?php echo $Lang['rescue_system']; ?></a>
              <?php endif; ?>
            </div>
          </div>
          <div class="btn-group ml-2 mr-2 mt-2">
            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i class="mdi mdi-chevron-down"></i></button>
            <div class="dropdown-menu">
              <?php if($Detail['dcim']['auth']['kvm'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="kvm">kvm</a>
              <?php endif; if($Detail['dcim']['auth']['ikvm'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="ikvm">ikvm</a>
              <?php endif; if($Detail['dcim']['auth']['novnc'] == 'on'): ?>
              <a class="dropdown-item dcim_service_module_button" href="javascript:void(0);"
                onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>');" data-func="vnc">vnc</a>
              <?php endif; ?>
            </div>
          </div>

          <?php elseif($Detail['module_button']['control']): ?>

          <div class="btn-group ml-2 mr-2 mt-2">

            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['control'] as $item): if($item['func'] != 'crack_pass' && $item['func'] != 'reinstall'): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

              <?php endif; ?>

              <?php endforeach; ?>

            </div>

          </div>

          <?php if($Detail['module_button']['console']): ?>

          <div class="btn-group ml-2 mr-2 mt-2">
            <?php if(($Detail['module_button']['console']|count) == 1): foreach($Detail['module_button']['console'] as $item): ?>
            <a class="btn btn-primary service_module_button d-flex align-items-center" href="javascript:void(0);"
              onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
              data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
            <?php endforeach; else: ?>
            <button type="button" class="btn btn-primary dropdown-toggle custom-button" data-toggle="dropdown"
              aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i class="mdi mdi-chevron-down"></i></button>

            <div class="dropdown-menu">

              <?php foreach($Detail['module_button']['console'] as $item): ?>

              <a class="dropdown-item service_module_button" href="javascript:void(0);"
                onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

              <?php endforeach; ?>

            </div>
            <?php endif; ?>

          </div>

          <?php endif; ?>

          <?php endif; ?>


          <!-- <div class="btn-group ml-2 mr-2 mt-2">
                <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>

                <span>

                  <?php if($Cancel['host_cancel']): ?>

                  <button class="btn btn-danger mb-1 h-100" id="cancelStopBtn" onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type'] ==
                    'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button>

                  <?php else: ?>

                  <button class="btn btn-primary mb-1 h-100" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>

                  <?php endif; ?>

                </span>

                <?php endif; ?>
              </div> -->
          <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.4-->
          <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
          <div class="btn-group ml-2 mr-2 mt-2">
            <span>
              <?php echo $item; ?>
            </span>
          </div>
          <?php endforeach; ?>
          <?php endif; ?>
          <!-- 结束 -->
        </div>

      </div>

    </div>

  </div>
  <div class="row">
    <div class="col-md-3">
      <div class="card">
        <div class="card-body">
          <!--
          <div class="mb-3 text-center" id="logininfo" style="width: 100px;height: 30px;line-height: 30px;background-color: #ffffff;box-shadow: 0px 4px 20px 2px rgba(6, 75, 179, 0.08);border-radius: 4px;cursor: pointer;">
            <?php echo $Lang['login_information']; ?>
          </div>
          -->

          <!-- 登录信息start -->
          <div class="row">
            <?php if($Detail['host_data']['domainstatus'] == 'Active'): ?>
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['login_information']; ?>
                  <i class="bx bx-user login-info-icon"></i>
                </p>
                <p class="mb-0"><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></p>
                <p class="mb-0"><?php echo $Lang['password']; ?>：
                  <!-- <span id="hidePwdBox" class="text-black">***********</span> -->
                  <?php if($Detail['host_data']['password'] == ''): ?>
                  <span class="text-black btnCopyPwd pointer dc">-</span>
                  <?php else: ?>
                  <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制"
                    id="copyPwdContent" class="text-black btnCopyPwd pointer dc"><?php echo $Detail['host_data']['password']; ?></span>
                  <?php endif; ?>
                  <!-- <i class="fas  fa-eye pointer ml-2 text-black" onclick="togglePwd()"></i> -->
                  <!-- <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPwd" data-clipboard-action="copy"
                    data-clipboard-target="#copyPwdContent"></i> -->
                  <?php if($Detail['host_data']['domainstatus']=="Active"): ?>
                  <span>
                    <button type="button"
                      class="btn btn-primary btn-sm waves-effect waves-light dcim_service_module_button fr rsb"
                      onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')" data-func="crack_pass"
                      data-type="default"><?php echo $Lang['crack_password']; ?></button>
                  </span>
                  <?php endif; ?>
                </p>
                <p class="mb-0"><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] ==
                  '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></p>

              </div>
            </div>
            <?php endif; if(($temp_custom_field_data = array_column($Detail['custom_field_data'], 'value', 'fieldname')) &&
            (isset($temp_custom_field_data['panel_address']) || isset($temp_custom_field_data['面板管理地址']) ||
            isset($temp_custom_field_data['panel_passwd']) || isset($temp_custom_field_data['面板管理密码']))): ?>
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray">
                <p class="text-gray">
                  <?php echo $Lang['panel_manage_info']; ?>
                  <i class="bx bx-receipt dc"></i>
                </p>
                <?php if(isset($temp_custom_field_data['panel_address']) || isset($temp_custom_field_data['面板管理地址'])): ?>
                <p class="mb-0"><?php echo $Lang['panel_manage_address']; ?>：<?php echo isset($temp_custom_field_data['panel_address']) ? $temp_custom_field_data['panel_address'] : $temp_custom_field_data['面板管理地址']; ?></p>
                <?php endif; if(isset($temp_custom_field_data['panel_passwd']) || isset($temp_custom_field_data['面板管理密码'])): ?>
                <!-- <p class="mb-0"><?php echo $Lang['panel_manage_password']; ?>：<span id="hidePanelPasswd">***********</span> -->
                <span data-toggle="popover" data-placement="top" data-trigger="hover" data-content="复制" id="panelPasswd"
                  class="btnCopyPanelPasswd dc pointer"><?php echo isset($temp_custom_field_data['panel_passwd']) ? $temp_custom_field_data['panel_passwd'] : ($temp_custom_field_data['面板管理密码'] ?:'--'); ?></span>
                <!-- <i class="fas fa-eye pointer ml-2 text-black" onclick="togglePanelPasswd()"></i>
                  <i class="bx bx-copy pointer ml-1 btn-copy text-black" id="btnCopyPanelPasswd" data-clipboard-action="copy" data-clipboard-target="#panelPasswd"></i> -->
                </p>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; ?>
          </div>
          <!-- 登录信息end -->

          <!-- 流量 -->
          <?php if(($Detail['host_data']['domainstatus'] == 'Active' || ($Detail['host_data']['domainstatus'] == 'Suspended' && $Detail['host_data']['suspendreason_type'] == 'flow')) && $Detail['host_data']['bwlimit'] > 0): ?>
          <!-- <div class="d-flex justify-content-end mb-2">
            <button type="button" class="btn btn-success btn-sm waves-effect waves-light"
              id="orderFlowBtn" onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
          </div> -->
          <div class="mt-4 mb-3">
            <i class="bx bx-circle" style="color:#f0ad4e"></i> <?php echo $Lang['used_flow']; ?>：<span id="usedFlowSpan">-</span>
            <i class="bx bx-circle" style="color:#34c38f"></i> <?php echo $Lang['residual_flow']; ?>：<span id="remainingFlow">-</span>
          </div>
          <div class="mb-4 ll-flex">
            <div class="progress w-75">
              <div class="progress-bar progress-bar-striped bg-success" id="totalProgress" role="progressbar"
                style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">100%
              </div>
            </div>
            <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb" id="orderFlowBtn"
              onclick="orderFlow($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['order_flow']; ?></button>
          </div>
          <?php endif; ?>
          <!-- 流量end -->
          <!-- 订购价格start -->
          <div class="row">
            <!-- <div class="col-12 my-2">
              <div class="d-flex justify-content-between align-items-center">
                <span><?php echo $Lang['first_order_price']; ?></span>
                <?php if($Detail['host_data']['status'] == 'Paid'): if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): ?>
                <button type="button" class="btn btn-primary btn-sm waves-effect waves-light" id="renew" onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['immediate_renewal']; ?></button>
                <?php endif; ?>
                <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                  <button type="button" class="btn btn-primary btn-sm waves-effect waves-light" id="renewpay"><?php echo $Lang['immediate_renewal']; ?></button>
                </a>
                <?php endif; ?>
              </div>
            </div> -->
            <div class="col-12 mb-2">
              <div class="bg-light card-body bg-gray pre-money-box">
                <p class="text-gray">
                  <?php echo $Lang['first_order_price']; ?>
                  <span class="fr"><?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] ==
                    'onetime'): else: ?><?php echo $Detail['host_data']['format_nextduedate']['msg']; ?><?php endif; ?></span>
                </p>
                <section class="d-flex align-items-center">
                  <h3 class="mb-0 mr-2 dc">
                    <?php echo !empty($Detail['host_data']['firstpaymentamount_desc']) ? $Detail['host_data']['firstpaymentamount_desc'] : '-'; ?></h3>

                  <!-- <span class="badge
                      <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?> - <?php else: ?><?php echo $Detail['host_data']['format_nextduedate']['msg']; ?><?php endif; ?></span> -->
                </section>

                <section class="d-flex align-items-center flex-wrap">
                  <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  $Detail['host_data']['status'] == 'Paid'): ?>
                  <span><?php echo $Lang['automatic_balance_renewal']; ?></span>
                  <div class="custom-control custom-switch custom-switch-md mb-4 ml-2" dir="ltr">
                    <input type="checkbox" class="custom-control-input" id="automaticRenewal"
                      onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew'] !=0): ?>checked
                      <?php endif; ?>> <label class="custom-control-label" for="automaticRenewal"></label>
                  </div>
                  <?php endif; if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                  ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): if($Detail['host_data']['status'] == 'Paid'): ?>
                  <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb" id="renew"
                    onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                  <?php endif; if($Detail['host_data']['status'] == 'Unpaid'): ?>
                  <a href="viewbilling?id=<?php echo $Detail['host_data']['invoice_id']; ?>">
                    <button type="button" class="btn btn-success btn-sm waves-effect waves-light rsb"
                      id="renewpay"><?php echo $Lang['renew']; ?></button>
                  </a>
                  <?php endif; ?>
                  <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): if($Cancel['host_cancel']): ?>
                  <!-- <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn"
                      onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php if($Cancel['host_cancel']['type']
                      ==
                      'Immediate'): ?><?php echo $Lang['stop_now']; else: ?><?php echo $Lang['stop_when_due']; ?><?php endif; ?></button> -->
                  <button class="btn btn-primary btn-sm rsb ml-2" id="cancelStopBtn" data-container="body"
                    data-toggle="popover" data-placement="top" data-trigger="hover"
                    data-content="将于{<?php echo date('Y-m-d',!is_numeric($Detail['host_data']['deletedate'])? strtotime($Detail['host_data']['deletedate']) : $Detail['host_data']['deletedate']); ?>}自动删除"
                    onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['cancel_out']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-danger btn-sm rsb ml-2" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['apply_out']; ?></button>
                  <?php endif; ?>
                  <?php endif; ?>
                </section>

                <section class="text-gray">
                  <p><?php echo $Lang['subscription_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></p>
                  <p><?php echo $Lang['payment_cycle']; ?>：<?php echo $Detail['host_data']['billingcycle_desc']; ?></p>

                  <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                  <p><?php echo $Lang['due_date']; ?>：-</p>
                  <?php else: ?>
                  <p><?php echo $Lang['due_date']; ?>：<?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></p>
                  <?php endif; ?>
                </section>
              </div>
            </div>
          </div>
          <!-- 订购价格end -->
          <!-- 配置项 -->
          <div class="row">
            <?php foreach($Detail['config_options'] as $item): if($item['option_type'] == '5'||$item['option_type'] == '6'||$item['option_type'] == '8'): else: ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['name']; ?>：<?php echo $item['sub_name']; ?>" class="bg-light card-body bg-gray mg-0 row">
                <p class="text-gray col-md-6 plr-0 text-nowrap">
                  <?php echo $item['name']; ?>
                </p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php if($item['option_type']===12): ?>
                  <img src="/upload/common/country/<?php echo $item['code']; ?>.png" width="20px">
                  <?php endif; ?>
                  <?php echo $item['sub_name']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>
            <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
            <div class="col-md-12 mb-2 configuration configuration_list">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $item['fieldname']; ?>：<?php echo $item['value']; ?>" class="bg-light card-body bg-gray mg-0 row">
                <p class="text-gray col-md-6 plr-0 text-nowrap"><?php echo $item['fieldname']; ?></p>
                <p class="mb-0 col-md-6 plr-0 text-nowrap text-right pl-2">
                  <?php echo $item['value']; ?>
                </p>
              </div>
            </div>
            <?php endif; ?>
            <?php endforeach; ?>
            <div onclick="isShowConfiguration()" class="configuration-btn-down isClick">查看更多信息</div>
            <div class="col-12 mb-2">
              <div data-toggle="popover" data-placement="top" data-trigger="hover"
                data-content="<?php echo $Lang['remarks_infors']; ?>：<?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>"
                class="bg-light card-body  bg-gray mg-0 row">
                <p class="text-gray col-md-3 plr-0 text-nowrap"><?php echo $Lang['remarks_infors']; ?></p>
                <p class="mb-0 col-md-9 plr-0 text-nowrap"><?php echo !empty($Detail['host_data']['remark']) ? $Detail['host_data']['remark'] : '-'; ?>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>

              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="col-md-9">
      <div class="card">
        <div class="card-body ">
          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">
            <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['dcimid'] && $Detail['dcim']['auth']['traffic'] ==
            'on'): ?>
            <li class="nav-item" id="chartLi">
              <a class="nav-link active" data-toggle="tab" href="#home1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="fas fa-home"></i></span> -->
                <span><?php echo $Lang['charts']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!($Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on')): ?>active<?php endif; ?>" data-toggle="tab"
                href="#profile1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="far fa-user"></i></span> -->
                <span><?php echo $Lang['upgrade_downgrade']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['download_data']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
              $Detail['host_data']['allow_upgrade_product']))): ?>active<?php endif; ?>" data-toggle="tab" href="#messages1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->
                <span><?php echo $Lang['file_download']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['show_traffic_usage']): ?>
            <li class="nav-item" id="usedLi">
              <a class="nav-link <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
              $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data']): ?>active<?php endif; ?>" data-toggle="tab"
                href="#dosage" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="far fa-envelope"></i></span> -->
                <span><?php echo $Lang['consumption']; ?></span>
              </a>
            </li>
            <?php endif; ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
              $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#settings1" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->
                <span><?php echo $Lang['journal']; ?></span>
              </a>
            </li>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#finance" role="tab">
                <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->
                <span><?php echo $Lang['finance']; ?></span>
              </a>
            </li>
          </ul>

          <!-- Tab panes -->
          <div class="tab-content p-3 text-muted">
            <?php if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on'): ?>
            <div class="tab-pane active" id="home1" role="tabpanel">
            </div>
            <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])): ?>
            <div class="tab-pane <?php if(!($Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on')): ?>active<?php endif; ?>" id="profile1"
              role="tabpanel">
              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-sm-3">

                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>
                        <div class="col-sm-6">

                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>
                        </div>
                        <div class="col-sm-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-sm-3">
                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                        </div>
                        <div class="col-sm-6">
                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>
                        </div>
                        <div class="col-sm-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; if($Detail['download_data']): ?>
            <div class="tab-pane <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product']))): ?>active<?php endif; ?>" id="messages1" role="tabpanel">
              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
            </div>
            <?php endif; ?>
            <div
              class="tab-pane <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data'] && !$Detail['host_data']['show_traffic_usage']): ?>active<?php endif; ?>"
              id="settings1" role="tabpanel">

            </div>

            <?php if($Detail['host_data']['show_traffic_usage']): ?>
            <div class="tab-pane <?php if(!($Detail['host_data']['domainstatus'] == 'Active' && $Detail['dcim'] && $Detail['dcim']['auth']['traffic'] == 'on') && !($Detail['host_data']['domainstatus'] == 'Active' && ($Detail['host_data']['allow_upgrade_config'] ||
            $Detail['host_data']['allow_upgrade_product'])) && !$Detail['download_data']): ?>active<?php endif; ?>" id="dosage"
              role="tabpanel">
              <div class="row d-flex align-items-center">
                <div class="col-md-3">
                  <input class="form-control" type="date" id="startingTime">
                </div>
                <span><?php echo $Lang['reach']; ?></span>
                <div class="col-md-3">
                  <input class="form-control" type="date" id="endTime">
                </div>
              </div>
              <div class="w-100 h-100">
                <div style="height: 500px" class="chart_content_box w-100" id="usedChartBox"></div>
              </div>
            </div>
            <?php endif; ?>
            <div class="tab-pane" id="finance" role="tabpanel">
              <!-- 财务 -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- 破解密码弹窗 -->
<div class="modal fade" id="dcimModuleResetPass" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0"><?php echo $Lang['crack_password']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <form id="crackPsdForm">
        <div class="modal-body">
          <div class="form-group row mb-0">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['password']; ?></label>
            <div class="col-md-6">
              <input type="text" class="form-control getCrackPsd" name="password">
            </div>
            <div class="col-md-1 fs-18 d-flex align-items-center">
              <i class="fas fa-dice create_random_pass pointer" onclick="create_random_pass()"></i>
            </div>
          </div>
          <label id="password-error-tip" class="control-label error-tip" for="password"></label>

          <div class="form-group row mb-4">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['crack_other']; ?></label>
            <div class="col-md-6">
              <div class="custom-control custom-checkbox">
                <input type="checkbox" class="custom-control-input" id="dcimModuleResetPassOther" value="1"
                  onchange="dcimModuleResetPassOther($(this))">
                <label class="custom-control-label" for="dcimModuleResetPassOther"><?php echo $Lang['password_will_cracked']; ?></label>
              </div>
            </div>
          </div>
          <div class="form-group row mb-4" style="display:none;" id="dcimModuleResetPassOtherUser">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['custom_user']; ?></label>
            <div class="col-md-6">
              <input type="text" class="form-control" name="user">
            </div>
          </div>
        </div>
        <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
      </form>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
        <button type="button" class="btn btn-primary submit"
          onclick="dcimModuleResetPass($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['determine']; ?></button>
      </div>
    </div>
  </div>
</div>

<!-- 救援系统弹窗 -->
<div class="modal fade" id="dcimModuleRescue" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0"><?php echo $Lang['rescue_system']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <form>
        <div class="modal-body">
          <div class="form-group row mb-4">
            <label for="horizontal-firstname-input"
              class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['system']; ?></label>
            <div class="col-md-8">
              <select class="form-control" name="system">
                <option value="1">Linux</option>
                <option value="2">Windows</option>
              </select>
            </div>
          </div>
        </div>
        <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
      </form>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
        <button type="button" class="btn btn-primary submit"
          onclick="dcimModuleRescue($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['determine']; ?></button>
      </div>
    </div>
  </div>
</div>

<!-- 重装系统弹窗 -->
<div class="modal fade" id="dcimModuleReinstall" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0"><?php echo $Lang['reinstall_system']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <form id="rebuildPsdForm">
        <div class="modal-body">
          <div class="row">
            <div class="col-md-2  d-flex align-items-center justify-content-end">
              <label class="float-right mb-0"><?php echo $Lang['system']; ?></label>
            </div>
            <div class="col-md-5">
              <div class="form-group mb-0">
                <select class="form-control configoption_os_group selectpicker" data-style="btn-default" name="os_group"
                  onchange="dcimModuleReinstallOsGroup($(this))">
                  <?php foreach($Detail['cloud_os_group'] as $item): if(strtolower($item['name'])=="windows"): $os_svg = '1'; elseif(strtolower($item['name'])=="centos"): $os_svg = '2'; elseif(strtolower($item['name'])=="ubuntu"): $os_svg = '3'; elseif(strtolower($item['name'])=="debian"): $os_svg = '4'; elseif(strtolower($item['name'])=="esxi"): $os_svg = '5'; elseif(strtolower($item['name'])=="xenserver"): $os_svg = '6'; elseif(strtolower($item['name'])=="freebsd"): $os_svg = '7'; elseif(strtolower($item['name'])=="fedora"): $os_svg = '8'; else: $os_svg = '9'; ?>
                  <?php endif; ?>
                  <option
                    data-content="<img class='mr-1' src='/upload/common/system/<?php echo $os_svg; ?>.svg' height='20'/><?php echo $item['name']; ?>"
                    value="<?php echo $item['id']; ?>"><?php echo $item['name']; ?></option>
                  <?php endforeach; ?>
                </select>
              </div>
            </div>
            <div class="col-md-5">
              <div class="form-group">
                <select class="form-control" name="os" data-os='<?php echo json_encode($Detail['cloud_os']); ?>'
                  onchange="dcimModuleReinstallOs($(this), '<?php echo $Detail['host_data']['os']; ?>')">
                  <?php foreach($Detail['cloud_os'] as $item): ?>
                  <option value="<?php echo $item['id']; ?>" data-group="<?php echo $item['group']; ?>"><?php echo $item['name']; ?></option>
                  <?php endforeach; ?>
                </select>
              </div>
            </div>
          </div>
          <div class="form-group row mb-0">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['password']; ?></label>
            <div class="col-md-6">
              <input type="text" class="form-control getRebuildPsd" name="password">
            </div>
            <div class="col-md-1 fs-18 d-flex align-items-center">
              <i class="fas fa-dice create_random_pass pointer" onclick="create_random_pass()"></i>
            </div>
          </div>
          <label id="password-error-tip-rebuild" class="control-label error-tip ml-8" for="password"></label>
          <div class="row mt-3">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['port']; ?></label>
            <div class="col-md-5">
              <input type="text" class="form-control" name="port" value="22">
            </div>
            <div class="col-md-1 fs-18 d-flex align-items-center">
              <i class="fas fa-dice module_reinstall_random_port"
                onclick="$('#dcimModuleReinstall input[name=\'port\']').val(parseInt(Math.random() * 65535))"></i>
            </div>
          </div>
          <div class="row" id="dcimModuleReinstallPart">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['partition_type']; ?></label>
            <div class="col-md-3">
              <div class="custom-control custom-radio">
                <input type="radio" class="custom-control-input" id="dcimModuleReinstallPart0" name="part_type"
                  onchange="showPartTypeConfirm('<?php echo $Detail['host_data']['os']; ?>')" value="0" checked="checked">
                <label class="custom-control-label" for="dcimModuleReinstallPart0"><?php echo $Lang['full_format']; ?></label>
              </div>
            </div>
            <div class="col-md-3">
              <div class="custom-control custom-radio">
                <input type="radio" class="custom-control-input" id="dcimModuleReinstallPart1" name="part_type"
                  onchange="showPartTypeConfirm('<?php echo $Detail['host_data']['os']; ?>')" value="1">
                <label class="custom-control-label"
                  for="dcimModuleReinstallPart1"><?php echo $Lang['first_partition_formatting']; ?></label>
              </div>
            </div>
            <div class="col-md-3">
              <div class="custom-control custom-checkbox">
                <input type="checkbox" class="custom-control-input" id="dcimModuleReinstallHigh"
                  onchange="showDcimDisk()">
                <label class="custom-control-label" for="dcimModuleReinstallHigh"><?php echo $Lang['senior']; ?></label>
              </div>
            </div>
          </div>
          <div class="row" id="dcimModuleReinstallDisk" style="display:none;">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['disk']; ?></label>
            <div class="col-md-6">
              <select class="form-control" name="disk">
                <option value="0">disk0</option>
              </select>
            </div>
          </div>
          <!--
          <div class="row" id="dcimModuleReinstallPartInfo">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end">分区</label>
            <div class="col-md-3">
              <div class="custom-control custom-radio">
                <input type="radio" class="custom-control-input" name="action" id="dcimModuleReinstallPartInfo0" onchange="reinstallPartInfoChange()"
                  value="0" checked="checked">
                <label class="custom-control-label" for="dcimModuleReinstallPartInfo0">默认</label>
              </div>
            </div>
            <div class="col-md-3">
              <div class="custom-control custom-radio">
                <input type="radio" class="custom-control-input" name="action" id="dcimModuleReinstallPartInfo1" onchange="reinstallPartInfoChange()"
                  value="1">
                <label class="custom-control-label" for="dcimModuleReinstallPartInfo1">附加配置</label>
              </div>
            </div>
          </div>
          <div class="row" style="display:none;" id="dcimModuleReinstallPartSetting">
            <label for="horizontal-firstname-input"
              class="col-md-2 col-form-label d-flex justify-content-end"></label>
            <div class="col-md-6">
              <select class="form-control" name="mcon">
                <option value="0">无可用附加分区配置</option>
              </select>
            </div>
          </div>
          -->
          <div class="row">
            <div class="col-md-2 offset-md-2 d-flex align-items-center justify-content-end">
              <div class="custom-control custom-checkbox mb-4">
                <input type="checkbox" class="custom-control-input" id="dcimModuleReinstallConfirm" value="1"
                  onchange="dcimReinstallConfirm($(this))">
                <label class="custom-control-label" for="dcimModuleReinstallConfirm"><?php echo $Lang['finished_backup']; ?></label>
              </div>
            </div>
          </div>
          <div class="row" id="dcimModuleReinstallPartMsg" style="display:none;">
            <div class="col-md-2"></div>
            <div class="col-md-9">
              <div class="part_error"></div>
            </div>
          </div>
          <div class="row">
            <div class="col-md-2"></div>
            <div class="col-md-9">
              <div id="dcimModuleReinstallMsg"></div>
            </div>
          </div>
        </div>
        <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
      </form>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
        <button type="button" class="btn btn-primary submit disabled" style="cursor:not-allowed;"
          onclick="dcimModuleReinstall($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['determine']; ?></button>
      </div>
    </div>
  </div>
</div>
<!-- 电源状态 -->
<script>
  var showPowerStatus = '<?php echo $Detail['module_power_status']; ?>';
  var powerStatus = {}

  $(function () {
    if (showPowerStatus == '1') {
      dcimGetPowerStatus('<?php echo app('request')->get('id'); ?>')
    }

    $('#powerBox').on('click', function () {
      dcimGetPowerStatus('<?php echo app('request')->get('id'); ?>')
    });
  })



  var timeOut = null
  var timeInterval = null
</script>

<script>
  // var showPWd = true
  $(function () {
    $(".nav-tabs-custom").find('.nav-item').find("a").removeClass('active')
    $(".nav-tabs-custom").find('.nav-item').eq(0).find("a").addClass('active')
    // 暂停状态悬浮原因
    $('.container-fluid').on('mouseover', '.status-suspended', function () {
      $('.xf-bg').show();
    })
    $('.container-fluid').on('mouseout', '.status-suspended', function () {
      $('.xf-bg').hide();
    })
    // 查看密码
    // $('#copyPwdContent').hide()
    // $('#hidePwdBox').hide()

    // 复制IP
    var clipboard = new ClipboardJS('.btn-copy-ip', {
      text: function (trigger) {
        return $('#copyIPContent').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    // 复制密码
    var clipboard = new ClipboardJS('.btnCopyPwd', {
      text: function (trigger) {
        return $('#copyPwdContent').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    // 一个ip时，不弹框，复制ip
    var clipboard = new ClipboardJS('.copyOneIp', {
      text: function (trigger) {
        return $('#copyOneIp').text()
      }
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })
  // function togglePwd() {
  //   showPWd = !showPWd

  //   if (showPWd) {
  //     $('#copyPwdContent').show()
  //     $('#hidePwdBox').hide()
  //   }
  //   if (!showPWd) {
  //     $('#copyPwdContent').hide()
  //     $('#hidePwdBox').show()
  //   }
  // }

  function togglePanelPasswd() {
    if ($("#hidePanelPasswd").is(':hidden')) {
      $("#hidePanelPasswd").show();
      $("#panelPasswd").hide();
    } else {
      $("#hidePanelPasswd").hide();
      $("#panelPasswd").show();
    }
  }
  //点击更多信息
  function isShowConfiguration(first = true) {
     let time = 0
    if (first) {
      time = 500
    } 
    for (let i = 0; i < 99; i++) {
      if (i > 3) {
        $('.configuration').eq(i).toggle(time)
      }
    }
    $('.configuration-btn-down').toggleClass('isClick')
    $('.configuration-btn-down').html('查看更多信息')
    $('.configuration-btn-down.isClick').html('收起更多信息')
    if($('.configuration_list').length < 4) {
      $('.configuration-btn-down').hide()
    }
  }
  isShowConfiguration(false)
  // 复制面板管理密码
  var clipboard_panelPasswd = new ClipboardJS('.btnCopyPanelPasswd', {
    text: function (trigger) {
      return $('#panelPasswd').text()
    }
  });
  clipboard_panelPasswd.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })
</script>

<script>
  // 图表tabs
  $(document).ready(function () {
    getComponentData()
  });

  let switch_id = []
  let chartsData = []
  let timeArray = []
  let name = []
  let typeArray = []
  let myChart = null

  $('#chartLi').on('click', function () {
    setTimeout(function () {
      myChart.resize()
    }, 0);
  });
  async function getChartDataFn(index) {
    selectTimeTypeFunc(index)
    const queryObj = {
      id: '<?php echo app('request')->get('id'); ?>',
      switch_id: switch_id[index],
      port_name: name[index],
      start_time: moment(timeArray[index].startTime).valueOf()
    }
    $.ajax({
      type: "post",
      url: '' + '/dcim/traffic',
      data: queryObj,
      success: function (data) {

        var obj = data.data.traffic || []
        var inArray = []
        var outArray = []
        var xName = []
        var inVal = []
        var outVal = []
        for (const item of obj) {
          if (item.type === 'in') {
            inArray.push(item)
            xName.push(moment(item.time).format('MM-DD HH:mm:ss'))
            inVal.push(item.value)
          } else if (item.type === 'out') {
            outArray.push(item)
            outVal.push(item.value)
          }
        }
        chartFunc(index, xName, inVal, outVal, data.data.unit)
        if (data.status === 200) {
          var obj = data.data.traffic
          var inArray = []
          var outArray = []
          var xName = []
          var inVal = []
          var outVal = []
          for (const item of obj) {
            if (item.type === 'in') {
              inArray.push(item)
              xName.push(moment(item.time).format('MM-DD HH:mm:ss'))
              inVal.push(item.value)
            } else if (item.type === 'out') {
              outArray.push(item)
              outVal.push(item.value)
            }
          }
          chartFunc(index, xName, inVal, outVal, data.data.unit)
        }
      }
    });
  }

  async function getComponentData() {
    const obj = {
      id: "<?php echo app('request')->get('id'); ?>"
    }
    $.ajax({
      type: "GET",
      url: '' + '/dcim/detail',
      data: obj,
      success: function (data) {
        if (data.status !== 200) {
          return
        }
        chartsData = data.data.switch ? data.data.switch : []
        let str = ``
        for (let i = 0; i < chartsData.length; i++) {
          const item = chartsData[i];
          timeArray.push({
            startTime: '',
            endTime: ''
          })
          typeArray.push({
            type: '7'
          })
          switch_id.push(chartsData[i].switch_id)
          name.push(chartsData[i].name)
          str += `<div
                    class="w-100 h-100">
                    <select class="form-control" id="chartSelect${i}" class="second_type" name="type" onchange="getChartDataFn(${i})">
                      <option value="24">24<?php echo $Lang['hour']; ?></option>
                      <option value="3">3<?php echo $Lang['day']; ?></option>
                      <option value="7" selected>7<?php echo $Lang['day']; ?></option>
                      <option value="30">30<?php echo $Lang['day']; ?></option>
                      <option value="999"><?php echo $Lang['whole']; ?></option>
                    </select>
                    <div style="height: 500px" class="w-100 h-100 d-flex justify-content-center" id="balanceCharts${i}"></div>
                    </div>`
        }
        $('#home1').append(str);
        for (let j = 0; j < chartsData.length; j++) {
          getChartDataFn(j)
        }
      }
    });
  }

  function selectTimeTypeFunc(index) {
    typeArray[index].type = $(`#chartSelect${index}`).val();

    if (typeArray[index].type === '7') { // 7天
      timeArray[index].startTime = moment(Date.now() - 7 * 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '3') { // 3天
      timeArray[index].startTime = moment(Date.now() - 3 * 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '30') { // 30天
      timeArray[index].startTime = moment(Date.now() - 30 * 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '24') { // 24h
      timeArray[index].startTime = moment(Date.now() - 24 * 3600 * 1000).format('YYYY-MM-DD')
    } else if (typeArray[index].type === '999') { // q全部
      timeArray[index].startTime = ''
    }
  }

  function chartFunc(index, xNameArray, inValArray, outValArray, unitY) {


    var inflow = '<?php echo $Lang['inflow_flow']; ?>'
    var outflow = '<?php echo $Lang['outflow_flow']; ?>'
    // 基于准备好的dom，初始化echarts实例
    myChart = echarts.init(document.getElementById('balanceCharts' + index))
    var option = {
      tooltip: {
        show: true,
        backgroundColor: '#fff',
        borderColor: '#eee',
        showContent: true,
        extraCssText: 'box-shadow: 0 1px 9px rgba(0, 0, 0, 0.1);',
        textStyle: {
          color: '#1e1e2d',
          textBorderWidth: 1
        },
        trigger: 'axis',
        axisPointer: {
          color: '#D9DAEA'
        }
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '10%',
        top: '8%',
        containLabel: true
      },
      color: ['#007bfc', '#3fbf70'],
      dataZoom: [ // 缩放
        {
          type: 'inside',
          throttle: 50
        }
      ],
      xAxis: {
        splitLine: {
          show: false
        },
        axisLine: {
          lineStyle: {
            color: '#1e1e2d'
          }
        },
        axisTick: {
          show: false
        },
        type: 'category',
        boundaryGap: false,
        data: xNameArray
      },
      yAxis: {
        axisLabel: {
          // formatter: '{value}' + unitY
          formatter: function (value) {
            return value + unitY
          }
        },
        axisLine: {
          show: false
        },
        minorTick: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#f5f4f8'
          }
        },
        type: 'value'
      },
      series: [{
          name: inflow,
          type: 'line',
          smooth: true,
          data: inValArray
        },
        {
          name: outflow,
          type: 'line',
          smooth: true,
          data: outValArray
        }
      ]
    }
    myChart.setOption(option, true) // true重绘
    window.addEventListener('resize', function () {
      myChart.resize()
    })
  }
</script>

<script>
  // 用量tabs
  let usedChart = null
  let usedStartTime
  let usedEndTime

  $(document).ready(function () {
    if ('<?php echo $Detail['host_data']['show_traffic_usage']; ?>') {
      chartOption()
    }
    if ($('#startingTime,#endTime').length > 0) getData()
    window.addEventListener('resize', function () {
      if (usedChart) usedChart.resize()
    })
  })
  $('#usedLi').on('click', function () {
    if (usedChart) {
      setTimeout(function () {
        usedChart.resize()
      }, 0);
    }
  });

  $('#startingTime,#endTime').on('change', function () {
    usedStartTime = $('#startingTime').val()
    usedEndTime = $('#endTime').val()
    getData()
  });
  // 图表配置
  function chartOption() {
    usedChart = echarts.init(document.getElementById('usedChartBox'))
    usedChart.setOption({
      backgroundColor: '#fff',
      title: {
        subtext: '',
        left: 'center',
        textAlign: 'left',
        subtextStyle: {
          lineHeight: 400
        }
      },
      tooltip: {
        backgroundColor: '#fff',
        padding: [10, 20, 10, 8],
        textStyle: {
          color: '#000',
          fontSize: 12
        },
        trigger: 'axis',
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#7dcb8f'
          }
        },
        formatter: function (params, ticket, callback) {
          // 
          const res = `
                    <div>
                      <div>` + '<?php echo $Lang['traffic_usage']; ?>' + `：${params[0].value}GB </div>
                      <div>${params[0].axisValue}</div>
                    </div>
                    `
          return res
        },
        extraCssText: 'box-shadow: 0px 4px 13px 1px rgba(1, 24, 167, 0.1);'
      },
      grid: {
        left: '80',
        top: 20,
        x: 70,
        x2: 50,
        y2: 80
      },
      xAxis: {
        offset: 15,
        type: 'category',
        data: [],
        boundaryGap: false,
        axisTick: {
          show: false
        },
        // 改变x轴颜色
        axisLine: {
          lineStyle: {
            type: 'dashed',
            color: '#ddd',
            width: 1
          }
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          }
        }
      },
      yAxis: {
        type: 'value',
        // 轴网格
        splitLine: {
          show: true,
          lineStyle: {
            color: '#ddd',
            type: 'dashed'
          }
        },
        axisTick: {
          show: false // 轴刻度不显示
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          show: true,
          textStyle: {
            color: '#999'
          },
          formatter: '{value}GB'
        }
      },
      series: [{
        name: '用量',
        type: 'line',
        smooth: true,
        showSymbol: true,
        symbol: 'circle',
        symbolSize: 3,
        // data: ['1200', '1400', '1008', '1411', '1026', '1288', '1300', '800', '1100', '1000', '1118', '123456'],
        data: [],
        areaStyle: {
          normal: {
            color: '#d4d1da',
            opacity: 0.2
          }
        },
        itemStyle: {
          normal: {
            color: '#0061ff' // 主要线条的颜色
          }
        },
        lineStyle: {
          normal: {
            width: 4,
            shadowColor: 'rgba(0,0,0,0.4)',
            shadowBlur: 10,
            shadowOffsetY: 10
          }
        }
      }]
    })
  }

  // 获取数据
  async function getData() {
    usedChart.showLoading({
      text: '<?php echo $Lang['data_loading']; ?>' + '...',
      color: '#999',
      textStyle: {
        fontSize: 30,
        color: '#444'
      },
      effectOption: {
        backgroundColor: 'rgba(0, 0, 0, 0)'
      }
    })



    const obj = {
      id: '<?php echo app('request')->get('id'); ?>',
      start: usedStartTime,
      end: usedEndTime
    }
    $.ajax({
      type: "get",
      url: '' + '/dcim/traffic_usage',
      data: obj,
      success: function (data) {

        usedChart.hideLoading()
        if (data.status !== 200) return false
        const xData = []
        const seriesData = [];
        (data.data || []).forEach(item => {
          xData.push(item.time)
          seriesData.push(item.value)
        })
        usedChart.setOption({
          title: {
            subtext: xData.length ? '' : '<?php echo $Lang['no_data_available']; ?>'
          },
          xAxis: {
            data: xData
          },
          series: [{
            data: seriesData
          }]
        })
        // 如果初始查询没有时间, 则设置默认时间为返回数据的第一个和最后一个时间
        if (!usedStartTime || !usedEndTime) {
          if (data.data.length) {
            usedStartTime = new Date().getFullYear() + '-' + data.data[0].time
            usedEndTime = new Date().getFullYear() + '-' + data.data[data.data.length - 1].time
            $('#startingTime').val(new Date().getFullYear() + '-' + data.data[0].time);
            $('#endTime').val(new Date().getFullYear() + '-' + data.data[data.data.length - 1].time);
          }
        }
      }
    });
  }

  // 分辨率改变, 重绘图表
  function resize() {
    usedChart.resize()
  }

  // 时间选择改变
  function dateChange() {
    const startTimeStamp = new Date(usedStartTime).getTime()
    const endTimeStamp = new Date(usedEndTime).getTime()
    if (usedStartTime && usedEndTime && startTimeStamp < endTimeStamp) {
      getData()
    }
  }
</script>

<script>
  // 获取基础数据
  const obj = {
    host_id: '<?php echo app('request')->get('id'); ?>'
  }
  $.ajax({
    type: "get",
    url: '' + '/host/dedicatedserver',
    data: obj,
    success: function (data) {
      const totalFlow = data.data.host_data.bwlimit // 总流量
      const usedFlow = data.data.host_data.bwusage // 已用流量
      const remainingFlow = (totalFlow - usedFlow).toFixed(1)
      let percentUsed = 100 - parseInt((usedFlow / totalFlow) * 100) || 0
      $('#totalProgress')
        .css('width', percentUsed + '%')
        .attr('aria-valuenow', percentUsed)
        .text(`${percentUsed}%`);

      $('#usedFlowSpan').text(`${usedFlow > 1024 ? ((usedFlow / 1024).toFixed(2) + 'TB') : (usedFlow + 'GB')}`);
      $('#remainingFlow').text(
        `${remainingFlow > 1024 ? ((remainingFlow / 1024).toFixed(2) + 'TB') : (remainingFlow + 'GB')}`);

      // 产品状态
      $('#statusBox').append(`<span class="sprite2 ${data.data.host_data.domainstatus}"></span>`)
    }
  });
</script>

<script>
  const logObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'log_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: logObj,
    success: function (data) {
      $(data).appendTo('#settings1');
    }
  });
  // 财务的append
  const financeObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'billing_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: financeObj,
    success: function (data) {
      $(data).appendTo('#finance');
    }
  });
</script>

<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<script>
  var getResintallStatusTimer = null;
  $(function () {
    getResintallStatus('<?php echo app('request')->get('id'); ?>');
  })
</script>

<script>
  var clipboard = null
  var clipboardpoppwd = null
  var ips = <?php echo json_encode($Detail['host_data']['assignedips']); ?>;
  // console.log('ips: ', ips);
  $(document).on('click', '#copyIPContent', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('IP地址')
    var iplist = ''
    if (clipboard) {
      clipboard.destroy()
    }
    for(let item in ips) {
      iplist += `
        <div>
          <span class="copyIPContent${item}">${ips[item]}</span>
          <i class="bx bx-copy pointer text-primary ml-1 btn-copy btnCopyIP${item}" data-clipboard-action="copy" data-clipboard-target=".copyIPContent${item}"></i>
        </div>
      `

      // 复制IP
      clipboard = new ClipboardJS('.btnCopyIP'+item, {
        text: function (trigger) {
          return $('.copyIPContent'+item).text()
        },
        container: document.getElementById('popModal')
      });
      clipboard.on('success', function (e) {
        toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
      })
    }

    $('#popContent').html(iplist)
  });


  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      <?php if($Detail['host_data']['domainstatus']=="Active"): ?>
	  <div>
      <button type="button" class="btn btn-primary btn-sm waves-effect waves-light dcim_service_module_button" onclick="dcim_service_module_button($(this), '<?php echo app('request')->get('id'); ?>')"
                  data-func="crack_pass" data-type="default"><?php echo $Lang['crack_password']; ?></button>
      </div>
	  <?php endif; ?>
    `)
  });


  $('#popModal').on('shown.bs.modal', function () {
    if (clipboardpoppwd) {
      clipboardpoppwd.destroy()
    }
    clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })

  // 破解密码
  $(document).on('blur', '.getCrackPsd', function () {
    veriCrackPsd()
  })

  function veriCrackPsd() {
    let result = checkingPwd($(".getCrackPsd").val(), passwordRules.num, passwordRules.upper, passwordRules.lower,
      passwordRules.special)
    if (result.flag) {
      $('#password-error-tip').css('display', 'none');
      $('.getCrackPsd').removeClass("is-invalid");
    } else {
      $("#password-error-tip").html(result.msg);
      $(".getCrackPsd").addClass("is-invalid");
      $('#password-error-tip').css('display', 'block');
    }
  }
  //重装系统
  $(document).on('blur', '.getRebuildPsd', function () {
    veriRebuildPsd()
  })

  function veriRebuildPsd() {
    let result = checkingPwd($(".getRebuildPsd").val(), passwordRules.num, passwordRules.upper, passwordRules.lower,
      passwordRules.special)
    if (result.flag) {
      $('#password-error-tip-rebuild').css('display', 'none');
      $('.getRebuildPsd').removeClass("is-invalid");
    } else {
      $("#password-error-tip-rebuild").html(result.msg);
      $(".getRebuildPsd").addClass("is-invalid");
      $('#password-error-tip-rebuild').css('display', 'block');
    }
  }

  $(function () {
    $("#crackPsdForm").on('click', ".create_random_pass", function (e) {
      veriCrackPsd()
    })
    $("#rebuildPsdForm").on('click', ".create_random_pass", function (e) {
      veriRebuildPsd()
    })
    getOsConf()
  })

  function getOsConf() {
    /*
    $.ajax({
      type: "POST",
      url: setting_web_url + '/provision/custom/<?php echo app('request')->get('id'); ?>',
      data: {
        func: 'getMirrorOsConfig',
      },
      success: function (res) {
        if(typeof res.data != 'undefined'){
          $('#dcimModuleReinstallPartSetting').data('data', res.data)
        }
      },
      error: function () {
        
      }
    })
    */
  }
</script>
<?php elseif($Detail['host_data']['type'] == "software"): ?>
    
<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<link rel="stylesheet" href="/themes/clientarea/default/assets/libs/bootstrap-select/css/bootstrap-select.min.css?v=<?php echo $Ver; ?>">
<script src="/themes/clientarea/default/assets/libs/bootstrap-select/js/bootstrap-select.min.js?v=<?php echo $Ver; ?>"></script>
<script type="text/javascript">
	$(function(){
		
	});
</script>
<script src="/themes/clientarea/default/assets/libs/moment/moment.js?v=<?php echo $Ver; ?>"></script>
<style>
  .btnzhuanyi a{
    color: #fff;
    background-color: #50a5f1;
    border-color: #50a5f1;
    padding: 4px 20px;
  }
  .btnzhuanyi a:hover{
     background-color: #50a5f1a6;
     border-color:#50a5f1a6;
  }
</style>
<div class="container-fluid">
  <div class="row">
    <div class="col-12 col-sm-8">
      <div class="card card_body p-3">
        <div class="row">
          <div class="col-12 col-sm-8 text-white">
            <div class="card card_body bg-primary p-4">
              <div class="mb-4">
                <span
                  class="status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?> px-2 py-0 rounded-sm fs-12">
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>
                </span>
              </div>
              <div class="mb-4"><?php echo $Detail['host_data']['productname']; ?></div>
              <div class="mb-4 text-white-50"><?php echo $Detail['host_data']['domain']; ?></div>
              <div class="mb-4">
                <ul class="pl-4">
                  <?php foreach($Detail['config_options'] as $configs): ?>
                  <li><?php echo $configs['name']; ?>: <?php echo $configs['sub_name']; ?></li>
                  <?php endforeach; foreach($Detail['custom_field_data'] as $fields): if($fields['showdetail'] == 1): ?>
                  <li><?php echo $fields['fieldname']; ?>: <?php echo $fields['value']; ?></li>
                  <?php endif; ?>
                  <?php endforeach; ?>
                </ul>
                <div>
                  <span class="text-white-50"><?php echo $Lang['remarks_infors']; ?>： <?php if($Detail['host_data']['remark']): ?><?php echo $Detail['host_data']['remark']; else: ?>-<?php endif; ?></span>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </div>
              </div>
              <?php if($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended'): ?>
              <div class="mb-4">
                <?php if($Detail['host_data']['cancel_control']): if($Cancel['host_cancel']): ?>
                  <button class="btn btn-danger mb-1" id="cancelStopBtn" onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['stop_when_due']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-info px-4 py-1 rounded-sm" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>
                  <?php endif; ?>
                <?php endif; ?>
                <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.5-->
                <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
                      <span class="btnzhuanyi">
                      <?php echo $item; ?>
                      </span>
                  <?php endforeach; ?>
                <?php endif; ?>
                <!-- 结束 -->
              </div>
              <?php endif; ?>
            </div>
          </div>
          <div class="col-12 col-sm-4 text-white">
            <div class="card card_body mb-3 p-3 bg-danger">
              <div>
                <i class="bx bx-shield"></i>
              </div>
              <div class="my-1"><?php echo $Lang['authorization_code']; ?></div>
              <div>
                <span id="copyCodeContent"><?php echo $Detail['host_data']['domain']; ?></span>
                <i class="bx bx-copy pointer text-white btn-copy" data-clipboard-action="copy"
                  data-clipboard-target="#copyCodeContent"></i>
              </div>
            </div>
            <div class="card card_body p-3 bg-warning">
              <div>
                <i class="bx bx-shield"></i>
              </div>
              <div class="my-1"><?php echo $Lang['ip_address']; ?></div>
              <div id="ipAddress"><?php echo !empty($Detail['host_data']['dedicatedip']) ? $Detail['host_data']['dedicatedip'] :  '-'; ?></div>
            </div>
          </div>
        </div>
        <div class="row">
          <div class="col-12 text-black-50">
            <div class="card card_body bg-light p-4">
              <div class="d-flex justify-content-between align-items-center mb-3">
                <div>
                  <i class="bx bx-bar-chart-alt"></i>
                  <span class="mr-2"><?php echo !empty($Detail['module_client_main_area'][0]) ? $Detail['module_client_main_area'][0]['name'] : $Lang['valid_domain_name']; ?></span>
                  <span id="validDomain"><?php echo !empty($Detail['module_client_main_area'][0]) ? $Detail['module_client_main_area'][0]['value'] :  '-'; ?></span>
                </div>
                <!-- <?php if($Detail['module_client_main_area'][0]): ?> -->
                <button id="resetBtn" onclick="resetAuth('<?php echo app('request')->get('id'); ?>')" class="btn btn-info px-3 py-1 rounded-sm"><?php echo $Lang['reset_authorization']; ?></button>
                <!-- <?php endif; ?> -->
              </div>
              <div class="d-flex justify-content-between align-items-center mb-3">
                <div>
                  <i class="bx bx-bar-chart-alt"></i>
                  <span class="mr-2"><?php echo !empty($Detail['module_client_main_area'][1]) ? $Detail['module_client_main_area'][1]['name'] : $Lang['valid_directory']; ?></span>
                  <span id="validPath"><?php echo !empty($Detail['module_client_main_area'][1]) ? $Detail['module_client_main_area'][1]['value'] :  '-'; ?></span>
                </div>
                <?php if($Detail['download_data']): ?>
                <a href="<?php echo $Detail['download_data']['0']['down_link']; ?>" class="bg-info px-3 py-1 rounded-sm text-white">下载文件</a>
                <?php endif; ?>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="col-12 col-sm-4">
      <div class="card card_body p-3">
        <div class="card card_body bg-light fs-12">
          <?php if($Renew['host']['billingcycle'] != 'free' && $Renew['host']['billingcycle'] != 'onetime'): ?>
          <div class="mt-4 ml-4">
            <?php if($Renew['host']['status'] == 'Unpaid'): ?>
            <span class="px-2 bg-danger rounded-sm fs-12 text-white"><?php echo $Lang['unpaid']; ?></span>
            <?php endif; if($Renew['host']['status'] == 'Paid'): ?>
            <span
              class="px-2 rounded-sm fs-12 text-white <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): ?> - <?php else: ?> <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span>
            <?php endif; ?>
          </div>
          <?php endif; ?>
          <div class="my-2 ml-4">
            <span
              class="font-weight-bold text-dark fs-18 mr-1"><?php echo $Renew['currency']['prefix']; ?><?php echo $Renew['host']['firstpaymentamount']; ?></span>
            <label for="" class="text-black-50 fz-12"><?php echo $Lang['first_order_price']; ?></label>
          </div>
          <ul class="text-black-50 fz-12">
            <li class="mb-2"><?php echo $Lang['ordering_time']; ?>： <?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></li>

            <li><?php echo $Lang['due_date']; ?>：
              <?php if($Renew['host']['status'] === 'Unpaid'): ?>-
              <?php else: if($Detail['host_data']['nextduedate'] && $Detail['host_data']['billingcycle'] != 'free' &&
                $Detail['host_data']['billingcycle'] != 'onetime'): ?>
                  <?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); else: ?>
                  <?php echo $Lang['not_due']; ?>
                <?php endif; ?>
              <?php endif; ?>
            </li>

          </ul>
        </div>
        <hr>
        <?php if($Renew['host']['billingcycle'] != 'free'): ?>
        <div class="font-weight-bold text-dark fs-14 mb-3"><?php echo $Renew['host']['status']=='Unpaid' ? $Lang['payment_information'] : $Lang['renewal_information']; ?></div>
        <div>
          <label class="text-black-50 fs-12 mr-2"><?php echo $Lang['pay_price']; ?></label>
          <span class="font-weight-bold text-dark fs-14"><?php echo $Renew['currency']['prefix']; ?><?php echo $Renew['host']['amount']; ?></span>
        </div>
        <div class="mb-2">
          <label class="text-black-50 fs-12 mr-2"><?php echo $Lang['payment_cycle']; ?></label>
          <?php if($Renew['host']['status'] == 'Unpaid'): ?>
          <span class="font-weight-bold text-dark fs-14"><?php echo $Renew['host']['billingcycle_zh']; ?></span>
          <?php endif; if($Renew['host']['status'] == 'Paid'): ?>
          <select class="form-control form-control-sm w-50 d-inline-block" name="" id="">
            <?php foreach($Renew['cycle'] as $cycle): ?>
            <option value="<?php echo $cycle['billingcycle']; ?>"><?php echo $cycle['billingcycle_zh']; ?></option>
            <?php endforeach; ?>
          </select>
          <?php endif; ?>
        </div>
        <?php if($Renew['host']['billingcycle'] != 'onetime' && $Renew['host']['status'] == 'Paid' && $Renew['host']['billingcycle'] != 'free'): ?>
        <div>
          <label class="text-black-50 fs-12 mr-2"><?php echo $Lang['automatic_balance_renewal']; ?></label>
          <div class="d-inline-block custom-control custom-switch custom-switch-md" dir="ltr">
            <input type="checkbox" class="custom-control-input" id="automaticRenewal" onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew'] !=0): ?>checked <?php endif; ?>>
            <label class="custom-control-label" for="automaticRenewal"></label>
          </div>
        </div>
        <?php endif; ?>
        <div>
          <?php if($Renew['host']['status'] == 'Paid'): ?>
          <span class="bg-primary px-3 py-1 pointer rounded-sm text-white" id="renew" onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['immediate_renewal']; ?></span>
          <?php endif; ?>
        </div>
        <?php endif; ?>
        <!-- end:非免费 -->
      </div>
    </div>
    <div class="col-12">
      <div class="card card_body p-4">
        <ul class="nav nav-tabs nav-tabs-custom" role="tablist">
          <li class="nav-item" onclick="getRechargeList('<?php echo app('request')->get('id'); ?>')">
            <a class="nav-link active" data-toggle="tab" href="#transaction" role="tab">              
              <span><?php echo $Lang['transaction_records']; ?></span>
            </a>
          </li>
          
          <?php if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
          <li class="nav-item">
            <a class="nav-link" data-toggle="tab" href="#downgrade" role="tab">
              <span><?php echo $Lang['upgrade_downgrade']; ?></span>
            </a>
          </li>
          <?php endif; if($Detail['download_data']): ?>
          <li class="nav-item">
            <a class="nav-link" data-toggle="tab" href="#filelist" role="tab">
              <span><?php echo $Lang['file_download']; ?></span>
            </a>
          </li>
          <?php endif; foreach($Detail['module_client_area'] as $item): ?>

          <li class="nav-item">

            <a class="nav-link" data-toggle="tab" href="#module_client_area_<?php echo $item['key']; ?>" role="tab">

              <!-- <span class="d-block d-sm-none"><i class="fas fa-cog"></i></span> -->

              <span><?php echo $item['name']; ?></span>

            </a>

          </li>

          <?php endforeach; ?>
          
          <li class="nav-item" onclick="getLogList('<?php echo app('request')->get('id'); ?>')">
            <a class="nav-link" data-toggle="tab" href="#log" role="tab">
              <span><?php echo $Lang['journal']; ?></span>
            </a>
          </li>
        </ul>
        <div class="tab-content p-3 text-white">
          <div class="tab-pane active" id="log" role="tabpanel">
            <div class="table-container">
              <div class="table-header mb-2">
                <div class="table-filter">
                  <div class="row">
                    <div class="col"></div>
                  </div>
                </div>
                <div class="table-search">
                  <div class="row justify-content-end">
                    <div class="col-sm-6">
                      <!-- <input type="text" class="form-control form-control-sm" id="logsearchInp" placeholder="<?php echo $Lang['search_by_keyword']; ?>"> -->
                    </div>
                  </div>
                </div>
              </div>
              <div class="table-responsive">
                <table class="table tablelist">
                  <colgroup>
                    <col>
                    <col>
                    <col>
                    <col>
                  </colgroup>
                  <thead class="bg-light">
                    <tr>
                      <th class="pointer" prop="create_time">
                        <span><?php echo $Lang['operation_time']; ?></span>
                        <span class="text-black-50 d-inline-flex flex-column justify-content-center ml-1 offset-3">
                          <i class="bx bx-caret-up"></i>
                          <i class="bx bx-caret-down"></i>
                        </span>
                      </th>
                      <th class="pointer" prop="detail">
                        <span><?php echo $Lang['operation_details']; ?></span>
                        <span class="text-black-50 d-inline-flex flex-column justify-content-center ml-1 offset-3">
                          <i class="bx bx-caret-up"></i>
                          <i class="bx bx-caret-down"></i>
                        </span>
                      </th>
                      <th class="pointer" prop="type">
                        <span><?php echo $Lang['operator']; ?></span>
                        <span class="text-black-50 d-inline-flex flex-column justify-content-center ml-1 offset-3">
                          <i class="bx bx-caret-up"></i>
                          <i class="bx bx-caret-down"></i>
                        </span>
                      </th>
                      <th class="pointer" prop="ip">
                        <span><?php echo $Lang['ip_address']; ?></span>
                        <span class="text-black-50 d-inline-flex flex-column justify-content-center ml-1 offset-3">
                          <i class="bx bx-caret-up"></i>
                          <i class="bx bx-caret-down"></i>
                        </span>
                      </th>
                    </tr>
                  </thead>
                  <tbody id="logTableData">
                    <!-- <tr>
                      <td>2021-02-01 18:00:00</td>
                      <td>取消重置密码成功 - Host ID:2641</td>
                      <td>智简魔方</td>
                      <td>106.81.231.193</td>
                    </tr> -->
                  </tbody>
                </table>
              </div>
              <!-- 表单底部调用开始 -->
              <!-- <div class="table-footer">
	<div class="table-tools">

	</div>
	<div class="table-pagination">
		<div class="table-pageinfo mr-2">
			<span><?php echo $Lang['common']; ?> <?php echo $Total; ?> <?php echo $Lang['strips']; ?></span>
			<span class="mx-2">
				<?php echo $Lang['each_page']; ?>
				<select name="" id="limitSel">
					<option value="10" <?php if($Limit==10): ?>selected<?php endif; ?>>10</option>
					<option value="15" <?php if($Limit==15): ?>selected<?php endif; ?>>15</option>
					<option value="20" <?php if($Limit==20): ?>selected<?php endif; ?>>20</option>
					<option value="50" <?php if($Limit==50): ?>selected<?php endif; ?>>50</option>
					<option value="100" <?php if($Limit==100): ?>selected<?php endif; ?>>100</option>
				</select>
				<?php echo $Lang['strips']; ?>
			</span>
		</div>
		<ul class="pagination pagination-sm">
			<?php echo $Pages; ?>
		</ul>
	</div>
</div>

<script>
	$(function () {

		// 每页数量选择改变
		$('#limitSel').on('change', function () {
			if ('<?php echo app('request')->get('action'); ?>') {
				location.href = '[url]?action=<?php echo app('request')->get('action'); ?>&keywords=<?php echo app('request')->get('keywords'); ?>&sort=<?php echo app('request')->get('sort'); ?>&orderby=<?php echo app('request')->get('orderby'); ?>&page=1&limit=' + $('#limitSel').val()
			} else {
				location.href = '[url]?keywords=<?php echo app('request')->get('keywords'); ?>&sort=<?php echo app('request')->get('sort'); ?>&orderby=<?php echo app('request')->get('orderby'); ?>&page=1&limit=' + $('#limitSel').val()
			}

		})


	})
</script> -->
            </div>
          </div>
          <div class="tab-pane" id="transaction" role="tabpanel">
            <div class="table-container">
              <div class="table-header mb-2">
                <div class="table-filter">
                  <div class="row">
                    <div class="col"></div>
                  </div>
                </div>
                <div class="table-search">
                  <div class="row justify-content-end">
                    <div class="col-sm-6">
                      <input type="text" class="form-control form-control-sm" id="transactionsearchInp"
                        placeholder="<?php echo $Lang['search_by_keyword']; ?>">
                    </div>
                  </div>
                </div>
              </div>
              <div class="table-responsive">
                <table class="table mb-0">
                  <thead class="thead-light">
                    <tr>
                      <th><?php echo $Lang['payment_time']; ?></th>
                      <th><?php echo $Lang['source']; ?></th>
                      <th><?php echo $Lang['payment_amount']; ?></th>
                      <th><?php echo $Lang['serial_number']; ?></th>
                      <th><?php echo $Lang['payment_method']; ?></th>
                    </tr>
                  </thead>
                  <tbody id="transactionTableData">
                    <!-- <tr>
                      <td>2021-02-01 18:00</td>
                      <td>来源</td>
                      <td>¥22.00</td>
                      <td>12564654654</td>
                      <td>微信支付</td>
                    </tr> -->
                  </tbody>
                </table>
              </div>
              <!-- 表单底部调用开始 -->
              <!-- <div class="table-footer">
	<div class="table-tools">

	</div>
	<div class="table-pagination">
		<div class="table-pageinfo mr-2">
			<span><?php echo $Lang['common']; ?> <?php echo $Total; ?> <?php echo $Lang['strips']; ?></span>
			<span class="mx-2">
				<?php echo $Lang['each_page']; ?>
				<select name="" id="limitSel">
					<option value="10" <?php if($Limit==10): ?>selected<?php endif; ?>>10</option>
					<option value="15" <?php if($Limit==15): ?>selected<?php endif; ?>>15</option>
					<option value="20" <?php if($Limit==20): ?>selected<?php endif; ?>>20</option>
					<option value="50" <?php if($Limit==50): ?>selected<?php endif; ?>>50</option>
					<option value="100" <?php if($Limit==100): ?>selected<?php endif; ?>>100</option>
				</select>
				<?php echo $Lang['strips']; ?>
			</span>
		</div>
		<ul class="pagination pagination-sm">
			<?php echo $Pages; ?>
		</ul>
	</div>
</div>

<script>
	$(function () {

		// 每页数量选择改变
		$('#limitSel').on('change', function () {
			if ('<?php echo app('request')->get('action'); ?>') {
				location.href = '[url]?action=<?php echo app('request')->get('action'); ?>&keywords=<?php echo app('request')->get('keywords'); ?>&sort=<?php echo app('request')->get('sort'); ?>&orderby=<?php echo app('request')->get('orderby'); ?>&page=1&limit=' + $('#limitSel').val()
			} else {
				location.href = '[url]?keywords=<?php echo app('request')->get('keywords'); ?>&sort=<?php echo app('request')->get('sort'); ?>&orderby=<?php echo app('request')->get('orderby'); ?>&page=1&limit=' + $('#limitSel').val()
			}

		})


	})
</script> -->
            </div>
          </div>
          <div class="tab-pane" id="filelist" role="tabpanel">
            <div class="table-container">
              <div class="table-responsive">
                <table class="table table-centered table-nowrap table-hover mb-0">
                  <thead>
                    <tr>
                      <th scope="col"><?php echo $Lang['file_name']; ?></th>
                      <th scope="col"><?php echo $Lang['upload_time']; ?></th>
                      <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
                    </tr>
                  </thead>
                  <tbody>
                    <?php if($Detail['download_data']): foreach($Detail['download_data'] as $item): ?>
                    <tr>
                      <td>
                        <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
                          <i class="mdi mdi-folder font-size-16 text-warning mr-2"></i>
                          <?php echo $item['title']; ?></a>
                      </td>
                      <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
                      <td><?php echo $item['downloads']; ?></td>
                      <td>
                        <div class="dropdown">
                          <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
                            <i class="bx bx-cloud-download"></i>
                          </a>
                        </div>
                      </td>
                    </tr>
                    <?php endforeach; else: ?>
                    <tr>
                      <td colspan="12">
                        <div class="no-data"><?php echo $Lang['nothing']; ?></div>
                      </td>
                    </tr>
                    <?php endif; ?>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
          <?php if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
          <div class="tab-pane" id="downgrade" role="tabpanel">
            <div class="container-fluid">
              <?php if($Detail['host_data']['allow_upgrade_product']): ?>
              <div class="row mb-3">
                <div class="col-12">
                  <div class="bg-light  rounded card-body">
                    <div class="row">
                      <div class="col-md-3">
                        <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                      </div>
                      <div class="col-md-6">
                        <span class="text-muted "><?php echo $Lang['upgrade_downgrade_two']; ?></span>
                      </div>
                      <div class="col-md-3">
                        <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                          id="upgradeProductBtn" onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
              <div class="row mb-3">
                <div class="col-12">
                  <div class="bg-light  rounded card-body">
                    <div class="row">
                      <div class="col-md-3">
                        <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                      </div>
                      <div class="col-md-6">
                        <span class="text-muted "><?php echo $Lang['upgrade_downgrade_description']; ?></span>
                      </div>
                      <div class="col-md-3">
                        <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                          id="upgradeConfigBtn" onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <?php endif; ?>
            </div>
          </div>
          <?php endif; foreach($Detail['module_client_area'] as $key=>$item): ?>
            <div class="tab-pane" role="tabpanel" id="module_client_area_<?php echo $item['key']; ?>">
              <div style="min-height: 550px;width:100%">
                <script>
                  $.ajax({
                    url : '/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>&date='+Date.parse(new Date()) 
                    ,type : 'get'
                    ,success : function(res) {
                        $('#module_client_area_<?php echo $item['key']; ?> > div').html(res);
                    }
                  })
                </script>
              </div>
              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                onload="this.height=$($('.main-content .card-body')[1]).height()-72" frameborder="0"
                width="100%"></iframe> -->
              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                frameborder="0" width="100%" style="min-height: 550px;"></iframe> -->

            </div>

          <?php endforeach; ?>
        </div>

      </div>
    </div>
  </div>
</div>

<script src="/themes/clientarea/default/assets/libs/clipboard/clipboard.min.js?v=<?php echo $Ver; ?>"></script>
<script>
  $(function () {
    getLogList('<?php echo app('request')->get('id'); ?>')

    if ('<?php echo $Detail['module_power_status']; ?>' == '1') {
      getNewStatus('<?php echo app('request')->get('id'); ?>')
    }
  })
  // 复制授权码
  var clipboard = new ClipboardJS('.btn-copy', {
    text: function (trigger) {
      return $('#copyCodeContent').text()
    }
  });
  clipboard.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

  $('#transactionsearchInp').on('keydown', function (e) {
    if (e.keyCode == 13) {
      getRechargeList('<?php echo app('request')->get('id'); ?>')
    }
  })

</script>
<?php elseif($Detail['host_data']['type'] == "cdn"): ?>
    
<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<div class="container-fluid">
  <div class="row">
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body">
          <div class="row">
            <div class="col-lg-5 mb-1">
              <div class="p-5 bg-primary rounded text-white d-flex flex-column
                                                justify-content-center align-items-center">
                <h1 class="text-white"><?php echo $Detail['host_data']['productname']; ?></h1>
                <p class="mb-4"><?php echo $Detail['host_data']['domain']; ?></p>
                <p>
                  <!-- 备注 -->
                  <span class="text-white-50"><?php echo $Lang['remarks_infors']; ?>： <?php if($Detail['host_data']['remark']): ?><?php echo $Detail['host_data']['remark']; else: ?>-<?php endif; ?></span>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>
                <span class="badge badge-pill py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>
                  mb-3">
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>
                </span>
              </div>
            </div>
            <div class="col-lg-7 mb-1">
              <div class="d-flex justify-content-between">
                <div class="table-responsive" style="min-height: auto;">
                  <table class="table mb-0 table-bordered">
                    <tbody>
                      <tr>
                        <th scope="row"><?php echo $Lang['price']; ?></th>
                        <td>
                          <?php echo $Detail['host_data']['firstpaymentamount_desc']; ?>
                          <span
                            class="ml-2 badge <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): else: ?> <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span>
                        </td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['subscription_date']; ?></th>
                        <td><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['payment_cycle']; ?></th>
                        <td><?php echo $Detail['host_data']['billingcycle_desc']; ?></td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['due_date']; ?></th>
                        <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                        <td>
                          <?php if($Detail['host_data']['billingcycle_desc'] == '一次性' || $Detail['host_data']['billingcycle_desc'] == '免费'): ?>
                            -
                          <?php else: ?>
                            <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?>
                          <?php endif; ?>
                        </td>
                        <?php else: ?>
                        <td><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></td>
                        <?php endif; ?>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['automatic_balance_renewal']; ?></th>
                        <td>
						<?php if($Detail['host_data']['billingcycle'] != 'onetime' && $Detail['host_data']['status'] == 'Paid' && $Detail['host_data']['billingcycle'] != 'free'): ?>
                          <div class="custom-control custom-switch custom-switch-md mb-3" dir="ltr">
                      <input type="checkbox" <?php if($Detail['host_data']['billingcycle_desc'] == '一次性' || $Detail['host_data']['billingcycle_desc'] == '免费'): ?> disabled <?php endif; ?> class="custom-control-input" id="automaticRenewal"
                              onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew']
                              !=0): ?>checked <?php endif; ?>> <label class="custom-control-label" for="automaticRenewal"></label>
                          </div>
						  <?php else: ?>
						  -
						  <?php endif; ?>
                        </td>
                      </tr>
                    </tbody>
                  </table>

                </div>

              </div>


              <div>
                
                <button type="button" class="btn btn-primary" id="logininfo">
                  <?php echo $Lang['login_information']; ?>

                  <i class="mdi mdi-chevron-down"></i>
                </button>

                <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): ?>
                <button type="button" class="btn btn-primary waves-effect waves-light" id="renew"
                  onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>
                <span>
                  <?php if($Cancel['host_cancel']): ?>
                  <button class="btn btn-danger" id="cancelStopBtn"
                    onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['stop_when_due']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-primary" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>
                  <?php endif; ?>
                </span>
                <?php endif; ?>
                <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.3-->
                <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
                    <div class="btn-group ml-0 mr-2">
                      <span>
                      <?php echo $item; ?>
                      </span>
                    </div>
                  <?php endforeach; ?>
                <?php endif; ?>
                <!-- 结束 -->

                <?php if($Detail['module_button']['control']): ?>

                <div class="btn-group">

                  <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i
                      class="mdi mdi-chevron-down"></i></button>

                  <div class="dropdown-menu">

                    <?php foreach($Detail['module_button']['control'] as $item): ?>

                    <a class="dropdown-item service_module_button" href="javascript:void(0);"
                      onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                      data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                      data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

                    <?php endforeach; ?>

                  </div>

                </div>
                <?php endif; if($Detail['module_button']['console']): ?>

                <div class="btn-group">
                  <?php if(($Detail['module_button']['console']|count) == 1): foreach($Detail['module_button']['console'] as $item): ?>
                  <a class="btn btn-primary service_module_button d-flex align-items-center" href="javascript:void(0);"
                    onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                    data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
                  <?php endforeach; else: ?>
                  <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i
                      class="mdi mdi-chevron-down"></i></button>

                  <div class="dropdown-menu">

                    <?php foreach($Detail['module_button']['console'] as $item): ?>

                    <a class="dropdown-item service_module_button" href="javascript:void(0);"
                      onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                      data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

                    <?php endforeach; ?>

                  </div>
                  <?php endif; ?>
                </div>
                <?php endif; ?>

              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body" style="min-height: 500px;">
          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">
            <?php foreach($Detail['module_client_area'] as $key=>$item): ?>

            <li class="nav-item">

              <a class="nav-link <?php if($key==0): ?>active<?php endif; ?>" data-toggle="tab"
                href="#module_client_area_<?php echo $item['key']; ?>" role="tab">
                <span><?php echo $item['name']; ?></span>

              </a>

            </li>
            <?php endforeach; if($Detail['config_options']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area']): ?>active<?php endif; ?>" data-toggle="tab" href="#profile1"
                role="tab">
                <span><?php echo $Lang['configuration_option']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area'] && !$Detail['config_options']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#downgrade" role="tab">
                <span><?php echo $Lang['upgrade_downgrade']; ?></span>
              </a>
            </li>
            <?php endif; ?>

            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area'] && !$Detail['config_options'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#finance" role="tab">
                <span><?php echo $Lang['financial_information']; ?></span>
              </a>
            </li>
            <?php if($Detail['download_data']): ?>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#download" role="tab">
                <span><?php echo $Lang['file_download']; ?></span>
              </a>
            </li>
            <?php endif; ?>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#settings1" role="tab">
                <span><?php echo $Lang['journal']; ?></span>
              </a>
            </li>

          </ul>

          <!-- Tab panes -->
          <div class="tab-content p-3 text-muted">
			<?php if($Detail['config_options']): ?>
            <div class="tab-pane  <?php if(!$Detail['module_client_area']): ?>active<?php endif; ?>" id="profile1" role="tabpanel">
              <div class="row">
                <?php foreach($Detail['config_options'] as $item): ?>
                <div class="col-md-2 mb-2">
                  <div class="bg-light">
                    <div class="card-body">
                      <p><?php echo $item['name']; ?></p>
                      <span><?php echo $item['sub_name']; ?></span>
                    </div>
                  </div>
                </div>
                <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
                <div class="col-md-2 mb-2">
                  <div class="bg-light">
                    <div class="card-body">
                      <p><?php echo $item['fieldname']; ?></p>
                      <span><?php echo $item['value']; ?></span>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
                <?php endforeach; ?>

              </div>
            </div>
			<?php endif; ?>
            <div class="tab-pane <?php if(!$Detail['module_client_area'] && !$Detail['config_options'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>" id="finance" role="tabpanel">

            </div>
            <div class="tab-pane" id="settings1" role="tabpanel">

            </div>
            <?php if($Detail['download_data']): ?>
            <div class="tab-pane" id="download" role="tabpanel">
              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
            </div>
            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <div class="tab-pane" id="downgrade" role="tabpanel">
              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-md-3">

                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>
                        <div class="col-md-6">
                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>
                        </div>
                        <div class="col-md-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-md-3">
                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                        </div>
                        <div class="col-md-6">
                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>
                        </div>
                        <div class="col-md-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; foreach($Detail['module_client_area'] as $key=>$item): ?>

            <div class="tab-pane <?php if($key==0): ?>active<?php endif; ?>" role="tabpanel"
              id="module_client_area_<?php echo $item['key']; ?>">

              <div class="width:100%">
                <script>
                  $.ajax({
                    url : '/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>&date='+Date.parse(new Date()) 
                    ,type : 'get'
                    ,success : function(res) {
                      $('#module_client_area_<?php echo $item['key']; ?> > div').html(res);
                    }
                  })
                </script>
              </div>

              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                onload="this.height=$($('.main-content .card-body')[1]).height()-72" frameborder="0"
                width="100%"></iframe> -->

            </div>

            <?php endforeach; ?>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="deactivateDia" style="display: none;">
    <form>
      <input type="hidden" value="<?php echo $Token; ?>" />
      <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>" />
      <div class="form-group row mb-4">
        <label class="col-sm-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
        <div class="col-sm-8">
          <select class="form-control" class="second_type" name="type">
            <option value="Immediate"><?php echo $Lang['remarks_infors']; ?>立即</option>
            <option value="Endofbilling" selected><?php echo $Lang['billing_cycle']; ?></option>
          </select>
        </div>
      </div>
      <div class="form-group row mb-0">
        <label class="col-sm-3 col-form-label text-right"><?php echo $Lang['cancelreason']; ?></label>
        <div class="col-sm-8">
          <div class="input-group">
            <select class="form-control" class="second_type" name="reason">
              <?php foreach($Detail['cancelist'] as $item): ?>
              <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
              <?php endforeach; ?>
            </select>
          </div>
        </div>
    </form>
  </div>
</div>

<style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<script src="/themes/clientarea/default/assets/libs/clipboard/clipboard.min.js?v=<?php echo $Ver; ?>"></script>
<script>
  function refresh(type) {
    location.reload();
  }


  // 查看密码
  var showPWd = false
  $('#copyPwdContent').hide()
  function togglePwd() {
    showPWd = !showPWd

    if (showPWd) {
      $('#copyPwdContent').show()
      $('#hidePwdBox').hide()
    }
    if (!showPWd) {
      $('#copyPwdContent').hide()
      $('#hidePwdBox').show()
    }
  }

  // 复制密码
  var clipboard = new ClipboardJS('.btn-copy', {
    text: function (trigger) {
      return $('#copyPwdContent').text()
    }
  });
  clipboard.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

</script>
<script>
  const logObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'log_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: logObj,
    success: function (data) {
      $(data).appendTo('#settings1');
    }
  });
  const financialObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'billing_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: financialObj,
    success: function (data) {
      $(data).appendTo('#finance');
    }
  });

  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      
    `)
  });


  $('#popModal').on('shown.bs.modal',function() {
    if (clipboardpoppwd) {
        clipboardpoppwd.destroy()
      }
     clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })
</script>
<?php elseif($Detail['host_data']['type'] == "other"): ?>
    
<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<div class="modal fade" id="popModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="popTitle"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body" id="popContent">
        
      </div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" data-dismiss="modal">确定</button>
			</div>
    </div>
  </div>
</div>
<div class="container-fluid">
  <div class="row">
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body">
          <div class="row">
            <div class="col-lg-5 mb-1">
              <div class="p-5 bg-primary rounded text-white d-flex flex-column
                                                justify-content-center align-items-center">
                <h1 class="text-white"><?php echo $Detail['host_data']['productname']; ?></h1>
                <p class="mb-4"><?php echo $Detail['host_data']['domain']; ?></p>
                <p>
                  <!-- 备注 -->
                  <span class="text-white-50"><?php echo $Lang['remarks_infors']; ?>： <?php if($Detail['host_data']['remark']): ?><?php echo $Detail['host_data']['remark']; else: ?>-<?php endif; ?></span>
                  <span class="bx bx-edit-alt pointer ml-2" data-toggle="modal" data-target="#modifyRemarkModal"></span>
                </p>
                <span class="badge badge-pill py-1 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>
                  mb-3">
                  <?php echo $Detail['host_data']['domainstatus_desc']; ?>
                </span>
              </div>
            </div>
            <div class="col-lg-7 mb-1">
              <div class="d-flex justify-content-between">
                <div class="table-responsive" style="min-height: auto;">
                  <table class="table mb-0 table-bordered">
                    <tbody>
                      <tr>
                        <th scope="row"><?php echo $Lang['price']; ?></th>
                        <td>
                          <?php echo $Detail['host_data']['firstpaymentamount_desc']; ?>
                          <span
                            class="ml-2 badge <?php echo $Detail['host_data']['format_nextduedate']['class']; ?>"><?php if($Detail['host_data']['format_nextduedate']['msg'] == '不到期'): else: ?> <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?> <?php endif; ?></span>
                        </td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['subscription_date']; ?></th>
                        <td><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['payment_cycle']; ?></th>
                        <td><?php echo $Detail['host_data']['billingcycle_desc']; ?></td>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['due_date']; ?></th>
                        <?php if($Detail['host_data']['billingcycle'] == 'free' || $Detail['host_data']['billingcycle'] == 'onetime'): ?>
                        <td>
                          <?php if($Detail['host_data']['billingcycle_desc'] == '一次性' || $Detail['host_data']['billingcycle_desc'] == '免费'): ?>
                            -
                          <?php else: ?>
                            <?php echo $Detail['host_data']['format_nextduedate']['msg']; ?>
                          <?php endif; ?>
                        </td>
                        <?php else: ?>
                        <td><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['nextduedate'])? strtotime($Detail['host_data']['nextduedate']) : $Detail['host_data']['nextduedate']); ?></td>
                        <?php endif; ?>
                      </tr>
                      <tr>
                        <th scope="row"><?php echo $Lang['automatic_balance_renewal']; ?></th>
                        <td>
						<?php if($Detail['host_data']['billingcycle'] != 'onetime' && $Detail['host_data']['status'] == 'Paid' && $Detail['host_data']['billingcycle'] != 'free'): ?>
                          <div class="custom-control custom-switch custom-switch-md mb-3" dir="ltr">
                      <input type="checkbox" <?php if($Detail['host_data']['billingcycle_desc'] == '一次性' || $Detail['host_data']['billingcycle_desc'] == '免费'): ?> disabled <?php endif; ?> class="custom-control-input" id="automaticRenewal"
                              onchange="automaticRenewal('<?php echo app('request')->get('id'); ?>')" <?php if($Detail['host_data']['initiative_renew']
                              !=0): ?>checked <?php endif; ?>> <label class="custom-control-label" for="automaticRenewal"></label>
                          </div>
						  <?php else: ?>
						  -
						  <?php endif; ?>
                        </td>
                      </tr>
                    </tbody>
                  </table>

                </div>

              </div>


              <div>
                
                <button type="button" class="btn btn-primary" id="logininfo">
                  <?php echo $Lang['login_information']; ?>

                  <i class="mdi mdi-chevron-down"></i>
                </button>

                <?php if($Detail['host_data']['billingcycle'] != 'free' && $Detail['host_data']['billingcycle'] != 'onetime' &&
                ($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Suspended')): ?>
                <button type="button" class="btn btn-primary waves-effect waves-light" id="renew"
                  onclick="renew($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['renew']; ?></button>
                <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' && $Detail['host_data']['cancel_control']): ?>
                <span>
                  <?php if($Cancel['host_cancel']): ?>
                  <button class="btn btn-danger" id="cancelStopBtn"
                    onclick="cancelStop('<?php echo $Cancel['host_cancel']['type']; ?>', '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['stop_when_due']; ?></button>
                  <?php else: ?>
                  <button class="btn btn-primary" data-toggle="modal"
                    data-target=".cancelrequire"><?php echo $Lang['out_service']; ?></button>
                  <?php endif; ?>
                </span>
                <?php endif; ?>
                <!--  20210331 增加产品转移hook输出按钮template_after_servicedetail_suspended.3-->
                <?php $hooks=hook('template_after_servicedetail_suspended',['hostid'=>$Detail['host_data']['id']]); if($hooks): foreach($hooks as $item): ?>
                    <div class="btn-group ml-0 mr-2">
                      <span>
                      <?php echo $item; ?>
                      </span>
                    </div>
                  <?php endforeach; ?>
                <?php endif; ?>
                <!-- 结束 -->

                <?php if($Detail['module_button']['control']): ?>

                <div class="btn-group">

                  <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false"><?php echo $Lang['control']; ?> <i
                      class="mdi mdi-chevron-down"></i></button>

                  <div class="dropdown-menu">

                    <?php foreach($Detail['module_button']['control'] as $item): ?>

                    <a class="dropdown-item service_module_button" href="javascript:void(0);"
                      onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                      data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>"
                      data-desc="<?php echo !empty($item['desc']) ? $item['desc'] : $item['name']; ?>"><?php echo $item['name']; ?></a>

                    <?php endforeach; ?>

                  </div>

                </div>
                <?php endif; if($Detail['module_button']['console']): ?>

                <div class="btn-group">
                  <?php if(($Detail['module_button']['console']|count) == 1): foreach($Detail['module_button']['console'] as $item): ?>
                  <a class="btn btn-primary service_module_button d-flex align-items-center" href="javascript:void(0);"
                    onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                    data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>
                  <?php endforeach; else: ?>
                  <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false"><?php echo $Lang['console']; ?> <i
                      class="mdi mdi-chevron-down"></i></button>

                  <div class="dropdown-menu">

                    <?php foreach($Detail['module_button']['console'] as $item): ?>

                    <a class="dropdown-item service_module_button" href="javascript:void(0);"
                      onclick="service_module_button($(this), '<?php echo app('request')->get('id'); ?>', '<?php echo $Detail['host_data']['type']; ?>')"
                      data-func="<?php echo $item['func']; ?>" data-type="<?php echo $item['type']; ?>" data-desc=""><?php echo $item['name']; ?></a>

                    <?php endforeach; ?>

                  </div>
                  <?php endif; ?>
                </div>
                <?php endif; ?>

              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body" style="min-height: 500px;">
          <ul class="nav nav-tabs nav-tabs-custom" role="tablist">
            <?php foreach($Detail['module_client_area'] as $key=>$item): ?>

            <li class="nav-item">

              <a class="nav-link <?php if($key==0): ?>active<?php endif; ?>" data-toggle="tab"
                href="#module_client_area_<?php echo $item['key']; ?>" role="tab">
                <span><?php echo $item['name']; ?></span>

              </a>

            </li>
            <?php endforeach; if($Detail['config_options']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area']): ?>active<?php endif; ?>" data-toggle="tab" href="#profile1"
                role="tab">
                <span><?php echo $Lang['configuration_option']; ?></span>
              </a>
            </li>
            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area'] && !$Detail['config_options']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#downgrade" role="tab">
                <span><?php echo $Lang['upgrade_downgrade']; ?></span>
              </a>
            </li>
            <?php endif; ?>

            <li class="nav-item">
              <a class="nav-link <?php if(!$Detail['module_client_area'] && !$Detail['config_options'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>"
                data-toggle="tab" href="#finance" role="tab">
                <span><?php echo $Lang['financial_information']; ?></span>
              </a>
            </li>
            <?php if($Detail['download_data']): ?>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#download" role="tab">
                <span><?php echo $Lang['file_download']; ?></span>
              </a>
            </li>
            <?php endif; ?>
            <li class="nav-item">
              <a class="nav-link" data-toggle="tab" href="#settings1" role="tab">
                <span><?php echo $Lang['journal']; ?></span>
              </a>
            </li>

          </ul>

          <!-- Tab panes -->
          <div class="tab-content p-3 text-muted">
			<?php if($Detail['config_options']): ?>
            <div class="tab-pane  <?php if(!$Detail['module_client_area']): ?>active<?php endif; ?>" id="profile1" role="tabpanel">
              <div class="row">
                <?php foreach($Detail['config_options'] as $item): ?>
                <div class="col-md-2 mb-2">
                  <div class="bg-light">
                    <div class="card-body">
                      <p><?php echo $item['name']; ?></p>
                      <span><?php echo $item['sub_name']; ?></span>
                    </div>
                  </div>
                </div>
                <?php endforeach; foreach($Detail['custom_field_data'] as $item): if($item['showdetail'] == 1): ?>
                <div class="col-md-2 mb-2">
                  <div class="bg-light">
                    <div class="card-body">
                      <p><?php echo $item['fieldname']; ?></p>
                      <span><?php echo $item['value']; ?></span>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
                <?php endforeach; ?>

              </div>
            </div>
			<?php endif; ?>
            <div class="tab-pane <?php if(!$Detail['module_client_area'] && !$Detail['config_options'] && !$Detail['host_data']['allow_upgrade_config'] && !$Detail['host_data']['allow_upgrade_product']): ?>active<?php endif; ?>" id="finance" role="tabpanel">

            </div>
            <div class="tab-pane" id="settings1" role="tabpanel">

            </div>
            <?php if($Detail['download_data']): ?>
            <div class="tab-pane" id="download" role="tabpanel">
              <div class="table-responsive">
  <table class="table table-centered table-nowrap table-hover mb-0">
    <thead>
      <tr>
        <th scope="col"><?php echo $Lang['file_name']; ?></th>
        <th scope="col"><?php echo $Lang['upload_time']; ?></th>
        <th scope="col" colspan="2"><?php echo $Lang['amount_downloads']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($Detail['download_data'] as $item): ?>
      <tr>
        <td>
          <a href="<?php echo $item['down_link']; ?>" class="text-dark font-weight-medium">
            <i
              class="<?php if($item['type'] == '1'): ?>mdi mdi-folder-zip text-warning<?php elseif($item['type'] == '2'): ?>mdi mdi-image text-success<?php elseif($item['type'] == '3'): ?>mdi mdi-text-box text-muted<?php endif; ?> font-size-16 mr-2"></i>
            <?php echo $item['title']; ?></a>
        </td>
        <td><?php echo date('Y-m-d H:i',!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['downloads']; ?></td>
        <td>
          <div class="dropdown">
            <a href="<?php echo $item['down_link']; ?>" class="font-size-16 text-primary">
              <i class="bx bx-cloud-download"></i>
            </a>
          </div>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
            </div>
            <?php endif; if($Detail['host_data']['allow_upgrade_config'] || $Detail['host_data']['allow_upgrade_product']): ?>
            <div class="tab-pane" id="downgrade" role="tabpanel">
              <div class="container-fluid">
                <?php if($Detail['host_data']['allow_upgrade_product']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-md-3">

                          <h5><?php echo $Lang['upgrade_downgrade']; ?></h5>
                        </div>
                        <div class="col-md-6">
                          <span><?php echo $Lang['upgrade_downgrade_two']; ?></span>
                        </div>
                        <div class="col-md-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeProductBtn"
                            onclick="upgradeProduct($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; if($Detail['host_data']['allow_upgrade_config']): ?>
                <div class="row mb-3">
                  <div class="col-12">
                    <div class="bg-light  rounded card-body">
                      <div class="row">
                        <div class="col-md-3">
                          <h5><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                        </div>
                        <div class="col-md-6">
                          <span><?php echo $Lang['upgrade_downgrade_description']; ?></span>
                        </div>
                        <div class="col-md-3">
                          <button type="button" class="btn btn-primary waves-effect waves-light float-right"
                            id="upgradeConfigBtn"
                            onclick="upgradeConfig($(this), '<?php echo app('request')->get('id'); ?>')"><?php echo $Lang['upgrade_downgrade_options']; ?></button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <?php endif; ?>
              </div>
            </div>
            <?php endif; foreach($Detail['module_client_area'] as $key=>$item): ?>

            <div class="tab-pane <?php if($key==0): ?>active<?php endif; ?>" role="tabpanel"
              id="module_client_area_<?php echo $item['key']; ?>">

              <div class="width:100%">
                <script>
                  $.ajax({
                    url : '/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>&date='+Date.parse(new Date()) 
                    ,type : 'get'
                    ,success : function(res) {
                      $('#module_client_area_<?php echo $item['key']; ?> > div').html(res);
                    }
                  })
                </script>
              </div>

              <!-- <iframe src="/provision/custom/content?id=<?php echo app('request')->get('id'); ?>&key=<?php echo $item['key']; ?>"
                onload="this.height=$($('.main-content .card-body')[1]).height()-72" frameborder="0"
                width="100%"></iframe> -->

            </div>

            <?php endforeach; ?>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="deactivateDia" style="display: none;">
    <form>
      <input type="hidden" value="<?php echo $Token; ?>" />
      <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>" />
      <div class="form-group row mb-4">
        <label class="col-sm-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
        <div class="col-sm-8">
          <select class="form-control" class="second_type" name="type">
            <option value="Immediate"><?php echo $Lang['remarks_infors']; ?>立即</option>
            <option value="Endofbilling" selected><?php echo $Lang['billing_cycle']; ?></option>
          </select>
        </div>
      </div>
      <div class="form-group row mb-0">
        <label class="col-sm-3 col-form-label text-right"><?php echo $Lang['cancelreason']; ?></label>
        <div class="col-sm-8">
          <div class="input-group">
            <select class="form-control" class="second_type" name="reason">
              <?php foreach($Detail['cancelist'] as $item): ?>
              <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
              <?php endforeach; ?>
            </select>
          </div>
        </div>
    </form>
  </div>
</div>

<style>
  .w-100{
    width: 100%;
  }
</style>
<div class="modal fade cancelrequire" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
  aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['out_service']; ?></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
          <input type="hidden" value="<?php echo $Token; ?>" />
          <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>" />

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['cancellation_time']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100"  name="type">
                <option value="Immediate"><?php echo $Lang['immediately']; ?></option>
                <option value="Endofbilling"><?php echo $Lang['cycle_end']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4">
            <label class="col-3 col-form-label text-right"><?php echo $Lang['reason_cancellation']; ?></label>
            <div class="col-8">
              <select class="form-control second_type w-100" name="temp_reason">
                <?php foreach($Cancel['cancelist'] as $item): ?>
                <option value="<?php echo $item['reason']; ?>"><?php echo $item['reason']; ?></option>
                <?php endforeach; ?>
                <option value="other"><?php echo $Lang['other']; ?></option>
              </select>
            </div>
          </div>

          <div class="form-group row mb-4" style="display:none;">
            <label class="col-3 col-form-label text-right"></label>
            <div class="col-8">
              <textarea class="form-control" maxlength="225" rows="3" placeholder="<?php echo $Lang['please_reason']; ?>" name="reason"
                value="<?php echo $Cancel['cancelist'][0]['reason']; ?>"></textarea>
            </div>
          </div>

        </form>

        <div class="modal-footer">
          <button type="button" class="btn btn-primary waves-effect waves-light" onClick="cancelrequest()"><?php echo $Lang['submit']; ?></button>
          <button type="button" class="btn btn-secondary waves-effect" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>

        </div>

      </div>
    </div>
  </div>
</div>



<script>

  var WebUrl = '/';
  $('.cancelrequire textarea[name="reason"]').val($('.cancelrequire select[name="temp_reason"]').val())
  $('.cancelrequire select[name="temp_reason"]').change(function () {
    if ($(this).val() == "other") {
      $('.cancelrequire textarea[name="reason"]').val('');
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').show();
    } else {
      $('.cancelrequire textarea[name="reason"]').val($(this).val())
      $('.cancelrequire textarea[name="reason"]').parents('.form-group').hide();
    }
  })

  function cancelrequest() {
    $('.cancelrequire').modal('hide');
    var content = '';
    var type = $('.cancelrequire select[name="type"]').val();
    if (type == 'Immediate') {
      content = '这将会立刻删除您的产品，操作不可逆，所有数据丢失';
    } else {
      content = '产品将会在到期当天被立刻删除，操作不可逆，所有数据丢失';
    }
    getModalConfirm(content, function () {
      $.ajax({
        url: WebUrl + 'host/cancel',
        type: 'POST',
        data: $('.cancelrequire form').serialize(),
        success: function (data) {
          if (data.status == '200') {
            toastr.success(data.msg);
            setTimeout(function () {
              window.location.reload();
            }, 1000)
          } else {
            toastr.error(data.msg);
          }
        }
      });
    })
  }
</script>
<script src="/themes/clientarea/default/assets/libs/clipboard/clipboard.min.js?v=<?php echo $Ver; ?>"></script>
<script>
  function refresh(type) {
    location.reload();
  }


  // 查看密码
  var showPWd = false
  $('#copyPwdContent').hide()
  function togglePwd() {
    showPWd = !showPWd

    if (showPWd) {
      $('#copyPwdContent').show()
      $('#hidePwdBox').hide()
    }
    if (!showPWd) {
      $('#copyPwdContent').hide()
      $('#hidePwdBox').show()
    }
  }

  // 复制密码
  var clipboard = new ClipboardJS('.btn-copy', {
    text: function (trigger) {
      return $('#copyPwdContent').text()
    }
  });
  clipboard.on('success', function (e) {
    toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
  })

</script>
<script>
  const logObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'log_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: logObj,
    success: function (data) {
      $(data).appendTo('#settings1');
    }
  });
  const financialObj = {
    id: '<?php echo app('request')->get('id'); ?>',
    action: 'billing_page'
  }
  $.ajax({
    type: "get",
    url: '' + '/servicedetail',
    data: financialObj,
    success: function (data) {
      $(data).appendTo('#finance');
    }
  });

  // 复制用户密码
  $(document).on('click', '#logininfo', function () {
    $('#popModal').modal('show')
    $('#popTitle').text('登录信息')

    $('#popContent').html(`
      <div><?php echo $Lang['user_name']; ?>：<?php echo $Detail['host_data']['username']; ?></div>
      <div>
        <?php echo $Lang['password']; ?>：<span id="poppwd"><?php echo $Detail['host_data']['password']; ?></span>
        <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="poppwdcopy" data-clipboard-action="copy" data-clipboard-target="#poppwd"></i>
      </div>
      <div><?php echo $Lang['port']; ?>：<?php if($Detail['host_data']['port'] == '0'): ?><?php echo $Lang['defaults']; else: ?><?php echo $Detail['host_data']['port']; ?><?php endif; ?></div>
      
    `)
  });


  $('#popModal').on('shown.bs.modal',function() {
    if (clipboardpoppwd) {
        clipboardpoppwd.destroy()
      }
     clipboardpoppwd = new ClipboardJS('#poppwdcopy', {
      text: function (trigger) {
        return $('#poppwd').text()
      },
      container: document.getElementById('popModal')
    });
    clipboardpoppwd.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  })
</script>
<?php elseif($Detail['host_data']['type'] == "ssl"): ?>
    
  <style>
    .seeInformationDiv {
      margin: 10px 0px;
    }

    .signContract-modal-dialog {
      max-width: 900px;
      margin: 1.75rem auto;
    }

    .contractfirstParty {
      display: flex;
      padding: 0px 70px;
    }

    .contractfirstPartyDiv div {
      height: 30px;
    }

    .contractfirstPartyDivTwo {
      margin-left: 100px;
    }

    .select_Div {
      display: flex;
      justify-content: flex-end;
    }



    .enclosure {
      padding: 0px 40px;
    }

    .cost {
      border-top: 1px solid #eff2f7;
      border-left: 1px solid #eff2f7;
      border-right: 1px solid #eff2f7;
      width: 100%;
      height: 40px;
      line-height: 40px;
    }

    #signature {
      width: 100%;
      display: flex;
      justify-content: center;
    }

    .signatureCanvas {
      border: 1px dashed #000 !important;
    }

    .display {
      display: none;
    }

    .afc-update-name {
      cursor: pointer;
      color: #1890ff;
    }



    #afcContent .control-label {
      margin-bottom: 0;
      text-align: right;
    }

    #recipientInfoContent .control-label {
      margin-bottom: 0;
      text-align: right;
    }

    .afc-modal-tip {
      color: #929292;
      margin-bottom: 10px;
    }

    @media screen and (max-width: 768px) {
      .contractfirstParty {
        display: block;
      }

      .contractfirstPartyDivTwo {
        margin-left: 0px;
      }

      .select_Div {
        display: block;
      }

      #afcContent .control-label {
        margin-bottom: 0;
        text-align: left;
      }
    }

    .defualt-add-tag {
      background: #b0d4f6;
      color: #fff;
      color: #fff !important;
      border-radius: 5px;
    }

    



    .ssl_p{
        font-size: 16px;font-weight: bold;color: #333333;padding-left: 12px;
    }

    .sslStatus{
        height: 116px;
        background: linear-gradient(270deg, #000000 0%, #F8F8FB 0%, #3D63FF 0%, #4D83FF 100%);
    }
    .sslStatusOne{
        display: flex;align-items: center;
        padding-left: 20px;
        justify-content: space-between;
    }
    .sslStatusOnediv{
      display: flex;align-items: center;
    }
    .sslStatusOnedivTwo{
      padding-right:50px;
    }
    .sslStatusOnediv>div:first-child>div{
        width: 55px;height: 53px;background: #F8F8FB;border-radius: 10px;
        display: flex;align-items: center;justify-content: center;
    }
    .sslStatusOne_status{
        margin: 0px  20px;
    }
    .sslStatusTwo{
        display: flex;
        align-items: center;
    }
    .sslStatusTwo>div:first-child{
        border-left:1px solid rgba(255, 255, 255, 0.25);
        border-right:1px solid rgba(255, 255, 255, 0.25);
        padding: 5px 0px 5px 20px;
        width: 230px;
    }
    .sslStatusTwo>div:last-child{
        border-right:1px solid rgba(255, 255, 255, 0.25);
        padding: 5px 0px 5px 20px;
        width: 230px;
    }
    .sslStatusTwoDiv>div:first-child{
        margin-bottom: 5px;
        font-size: 14px;
        color: #A3BCFF;
    }
    .sslStatusTwoDiv>div:last-child{
        font-size: 16px;
        color: #FFFFFF;
    }
    .sslStatusThree{
        padding-right: 20px;
    }
    .sslStatusThree>.btn{
        background: #6F87FC;
        box-shadow: 2px 2px 9px rgba(0, 28, 144, 0.15);
        opacity: 1; 
        border-radius: 4px;
        color: #fff;
        border: 1px solid #6F87FC;
        padding-left:15px;
        padding-right:15px;
    }
    .ssl_operations{
        text-align: right;
    }
    @media screen and (max-width: 768px) {
      .sslStatus {
        height: 100% !important;
      }
      .sslStatusOne{
          padding: 0px;
      }
      .sslStatusOnedivTwo{
        padding-right:0px !important;
      }
      .ssl_operations{
          text-align: left !important;
      }
      .ssl_domain{
          margin: 15px 0px;
      }
      .ssl_detailDiv>div{
        padding-left: 0px !important;
    }
      .d-j-a-center-operation{
        display: block !important;
      }
      .td_width{
        /* max-width: auto !important; */
        /* min-width: auto !important; */
        width:auto !important;
      }
      .label-text-right{
        text-align: left !important;
      }
      .sslIssueDiv{
        height: auto !important;
      }
    }
    .ssl_orderDetail{
        display: flex;
    }
    .ssl_orderDetail>label{
        color: #999999;
        max-width: 100px;
        min-width: 100px;
        text-align: right;
    }
    .ssl_detailDiv>div{
        padding-left: 20px;
    }
    
    .stepDiv{
        margin:30px 15px;
        
    }
    .step{
        padding: 5px 25px 30px;
        border-left: 1px solid #EFEFFF;
        position: relative;
    }
    .step>.yuan{
        width: 30px;
        height: 30px;
        border-radius: 50%;
        text-align: center;
        line-height: 28px;
        left: -15px;
        border:1px solid;
        top: 0px;
        position: absolute;
    }
    .stepContent_text{
        margin:10px 0px;
    }
    .stepNoborder{
        border-color: #fff;
    }
    .ssl-dialog-tip>div{
        color: #006EFF;
        height: 32px;
        display: flex;
        align-items: center;
        background: #F3F8FF;
        opacity: 1;
    }
    .ssl-more-operations li a{
      color: #212529;
      padding: .35rem 0.5rem;
    }
    .d-j-a-center {
      display: flex;
      justify-content: center;
      align-items: center;
      flex-flow: wrap;
    }
    .verification-method{
      display: flex;
      width: 60%;
    }
    .verification-method>div{
      width:50%;height: 34px;border:1px solid #D1D6DD;
      text-align: center;
      line-height: 32px;
      font-size: 12px;
    }
    .verification-method>div:first-child{
      border-radius: 3px 0px 0px 3px;
    }
    .verification-method>div:last-child{
      border-radius: 0px 3px 3px 0px;
      border-left: 0px;
    }
    .poitner {
      cursor: pointer;
    }
    .v-active{
      background: #6064FF;
      color: #fff;
    }
    .d-j-a-center-operation{
      display: flex;
    }
    #v-operation-text{
      margin-bottom: 20px;
      background: #eff2f7;
      padding: 5px 0px;
      font-weight: bold;
    }
    #v-operation-text-two{
      margin-bottom: 20px;
      background: #eff2f7;
      padding: 5px 0px;
      font-weight: bold;
    }
    .td_width{
      /* max-width: 200px; */
      /* min-width: 200px; */
      width:200px;
      text-align: center;
    }
    #file-table{
      display: none;
    }
    .issueCsrDiv{
      display: flex;
      align-items: center;
    }
    .issueCsrDiv >div{
      display: flex;
      align-items: center;
    }
    .issueCsrDiv >div>span{
      margin-left: 10px;
    }
    .issueCsrDiv >div>input{
      width: 16px;
      height: 16px;
    }
    .issueCsrDiv div:last-child{
      margin-left: 20px;
    }
    .issueContactsDiv{
      display: flex;
      align-items: center;
    }
    .issueContactsDiv>input{
      width: 16px;
      height: 16px;
    }
    .verificationDiv .d-j-a-center>label:before{
      content: '*';
      color: red;
    }
    .redStar:before{
      content: '*';
      color: red;
    }
    .red_label:before{
      content: '*';
      color: red;
    }
    .label-text-right{
      text-align: right;
    }
    .sslIssueDiv{
      overflow-y: auto;
      height: 450px;
      overflow-x: hidden;
    }
    /* 未完成步骤 */
    .sslIncompleteStep > .yuan {
      border-color: #999999;
      background: #fff;
      color: #999999;
    }
    .sslIncompleteStep > .stepContent {
      color: #999999;
    }
    .sslIncompleteStep  .btn-ssl-b{
      color: #999999;
      background-color: #fff;
      border-color: #EFEFFF;
    }
    .sslIncompleteStep  .btn-ssl{
      color: #999999;
      background-color: #fff;
      border-color: #EFEFFF;
    }
    /* 当前步骤 */
    .sslCurrentStep > .yuan{
      background: #6064FF;
      border-color:#6064FF; 
      color: #fff;
    }
    .sslCurrentStep  .btn-ssl-b{
      color: #fff;
      background-color: #6064ff;
      border-color: #6064ff;
    }
    .sslCurrentStep  .btn-ssl{
      border:1px solid #6064FF;color: #6064FF;
    }
    /* 已完成步骤 */
    .sslCompletedStep  > .yuan {
      background: #DFE0FF;
      border-color: #DFE0FF;
      color: #fff;
    }
    .sslCompletedStep  >.stepContent{
      color: #999;
    }

    #showDomain{
      color:#426dff;display: flex;align-items: center;
    }
    #showDomain span:first-child {
      display: inline-block;font-size: 24px;width: 26px;height: 26px;border: 1px solid #426dff;line-height: 22px;text-align: center;border-radius: 50%;
    }
    #showDomain span:last-child{
      margin-left: 5px;
    }
    #hideDomain{
      display: none;
      color:#999999;align-items: center;
    }
    #hideDomain span:first-child{
      display: inline-block;font-size: 24px;width: 26px;height: 26px;border: 1px solid #999999;line-height: 22px;text-align: center;border-radius: 50%;
    }
    #hideDomain span:last-child{
      margin-left: 5px;
    }

    /* ssl证书待验证 */
    .status-verifiy_active{
      color: #fff;
      background-color: #fca426 !important;
      border-color: #fca426 !important;
      padding: 5px 10px !important;
    }
    /* ssl证书即将过期 */
    .status-overdue_active{
      color: #fff;
      background-color: #6064FF !important;
      border-color: #6064FF !important;
      padding: 5px 10px !important;
    }
    /* ssl证书已签发 */
    .status-issue_active{
      color: #fff;
      background-color: #3FBF70 !important;
      border-color: #3FBF70 !important;
      padding: 5px 10px !important;
    }
    /* ssl证书未付款 */
    .status-pending{
      color: #fff;
      background-color: #EE6161 !important;
      border-color: #EE6161 !important;
      padding: 5px 10px !important;
    }
    /* ssl证书未使用 */
    .status-active{
      color: #fff;
      background-color: #67A4FF !important;
      border-color: #67A4FF !important;
      padding: 5px 10px !important;
    }
    .borderRed{
      border-color:red;
    }
    #ipdomainNames:hover {
    color: #FCA426
  }
  </style>


<style>
  .paySelect {
    float: right;
    width: 25%;
    height: 1.5rem;
    border-color: #ddd;
    margin-left: 42%;
  }

  @media screen and (max-width: 440px) {
    .paySelect {
      width: 33%;
    }

  }

  @media screen and (max-width: 375px) {
    .paySelect {
      margin-left: 21%;
    }
  }

  @media screen and (max-width: 320px) {
    .paySelect {
      margin-left: 7%;
      width: 36%;
    }
  }
</style>
<script>
  var _url = '';

  function payTypeChange(invoiceid) {
    let paymt = 1;
    var url = _url + '/change_paymt';
    if ($('.paySelect option:selected').val() == 0) {
      paymt = 0
    }
    let invoiceidNew = $("#myLargeModalLabel").html().split(' - ')[1]

    $.ajax({
      type: "POST",
      data: {
        invoiceid: invoiceidNew,
        paymt: paymt
      },
      url: url,
      success: function (data) {
        if ($('.paySelect option:selected').val() == 0) {
          payamount(invoiceidNew, 0)
        } else {
          payamount(invoiceidNew)
        }
      }
    })
  }

  $(function () {
    $(".close").click(function () {
      $(".modal-header .paySelect").remove();
    })
  })
</script>
<div id="pay">
  <div class="modal fade pay" id="myModal" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
    aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered ">
      <div class="modal-content">
        <div class="modal-header">
          <h6 class="modal-title mt-0" id="myLargeModalLabel"></h6>
          <!--<?php if($Action != 'recharge'): ?>
                <select class="paySelect" onchange="payTypeChange(<?php echo $Pay['invoiceid']; ?>)">
                  <option value="0">现金支付</option>
                    <?php if((!empty($paymt['is_open_credit_limit']) && $paymt['credit_limit_balance'] >= $paymt['subtotal'])): ?>
                        <option value="1">信用额支付</option>
                    <?php endif; ?>
                </select>
                <?php endif; ?>-->
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">

        </div>
      </div>
    </div>
  </div>

</div>


<div id="loading-icon" style="display:none">
  <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
    style="margin:auto;background:#fff;display:block;" width="200px" height="200px" viewBox="0 0 100 100"
    preserveAspectRatio="xMidYMid">
    <g>
      <circle cx="73.801" cy="68.263" fill="#93dbe9" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="0s">
        </animateTransform>
      </circle>
      <circle cx="68.263" cy="73.801" fill="#689cc5" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.062s">
        </animateTransform>
      </circle>
      <circle cx="61.481" cy="77.716" fill="#5e6fa3" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.125s">
        </animateTransform>
      </circle>
      <circle cx="53.916" cy="79.743" fill="#3b4368" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.187s">
        </animateTransform>
      </circle>
      <circle cx="46.084" cy="79.743" fill="#191d3a" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.25s">
        </animateTransform>
      </circle>
      <circle cx="38.519" cy="77.716" fill="#d9dbee" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.312s">
        </animateTransform>
      </circle>
      <circle cx="31.737" cy="73.801" fill="#b3b7e2" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.375s">
        </animateTransform>
      </circle>
      <circle cx="26.199" cy="68.263" fill="#93dbe9" r="3">
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50"
          times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.437s">
        </animateTransform>
      </circle>
      <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;0 50 50" times="0;1"
        keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s">
      </animateTransform>
    </g>
  </svg>
</div>
<body>
  <div class="card mb-4 bg-primary sslStatus mb-4">
    <div class="card-body">
        <div class="row align-items-center text-white" style="height: 100%;">
            <div class="col-sm-12 col-md-4">
                <div class="sslStatusOne">
                  <div class="sslStatusOnediv">
                    <div>
                        <div>
                            <img src="/static/images/<?php echo $Detail['host_data']['domainstatus']; ?>.png" alt="">
                        </div>
                    </div>
                    <div class="sslStatusOne_status">
                        <span style="font-size: 20px;font-weight: bold;"><?php echo $Detail['host_data']['productname']; ?></span>
                        <span class="badge badge-pill font-size-12 status-<?php echo strtolower($Detail['host_data']['domainstatus']); ?>"><?php echo $Detail['host_data']['domainstatus_desc']; ?></span>
                    </div>
                  </div>
                  <div class="sslStatusOnedivTwo">
                      <?php if($Detail['host_data']['domainstatus'] == 'Pending'): ?>
                          <a href="#">
                              <button type="button" style="width:80px;" class="btn btn-primary btn-sm waves-effect waves-light" onclick="payamount(<?php echo $Detail['host_data']['invoice_id']; ?>);">去付款</button>
                          </a>
                      <?php endif; ?>
                      <!-- 未使用 ->签发 -->
                      <?php if($Detail['host_data']['domainstatus'] == 'Active'): ?>
                          <a href="#">
                              <button type="button" style="width:80px;" id="IssueBtn" class="btn btn-primary btn-sm waves-effect waves-light issueCertBtn">签发</button>
                          </a>
                      <?php endif; if(($Detail['host_data']['domainstatus'] != 'Pending' && $Detail['host_data']['domainstatus'] != 'Active' && $Detail['host_data']['domainstatus'] != 'Deleted')): ?>
                          <a href="#">
                              <button type="button" style="min-width:80px;" id="againIssueBtn" class="btn btn-primary btn-sm waves-effect waves-light issueCertBtn">重新签发</button>
                          </a>
                      <?php endif; ?>


                  </div>
                </div>
            </div>
            <div class="col-sm-12 col-md-4 ssl_domain">
                <div class="sslStatusTwo">
                    <div class="sslStatusTwoDiv">
                        <div>域名</div>
                        <div class="poitner" id="ipdomainNames" <?php if($Detail['host_data']['certssl_orderinfo']['domainNames']): ?> data-toggle="modal" data-target="#sslIPdialog" <?php endif; ?>>
                            <?php echo isset($Detail['host_data']['certssl_orderinfo']['used_domainNames']) ? $Detail['host_data']['certssl_orderinfo']['used_domainNames'] : '-'; if(($Detail['host_data']['certssl_orderinfo']['domainNames'])): ?>
                                (<?php echo count($Detail['host_data']['certssl_orderinfo']['domainNames_arr']); ?>)
                            <?php endif; ?>
                        </div>
                    </div>
                    <div class="sslStatusTwoDiv">

                        <?php if(($Detail['host_data']['domainstatus'] != 'Issue_Active' && $Detail['host_data']['domainstatus'] != 'Deleted' && $Detail['host_data']['domainstatus'] != 'Overdue_Active')): ?>
                            <div>购买周期</div>
                            <div>
                                <?php foreach($Detail['config_options'] as $k => $val): if($val['name_k'] == 'certssl_cert_year'): ?>
                                        <?php echo $val['sub_name']; ?>
                                    <?php endif; ?>
                                <?php endforeach; ?>
                            </div>
                        <?php else: ?>
                            <div>到期时间</div>
                            <div>
                                <?php if(($Detail['host_data']['domainstatus'] == 'Deleted')): ?>
                                    已过期
                                <?php endif; if(($Detail['host_data']['domainstatus'] == 'Issue_Active' || $Detail['host_data']['domainstatus'] == 'Overdue_Active')): ?>
                                    <?php echo $Detail['host_data']['certssl_orderinfo']['due_time_day']; ?>天后到期
                                <?php endif; ?>
                            </div>
                        <?php endif; ?>
                    </div>
                </div>
            </div>
            <!-- <div class="col-sm-12 col-md-4 ssl_operations" style="text-align: right;">
                <div class="btn-group sslStatusThree">
                    <button type="button" class="btn btn-default dropdown-toggle btn-sm" data-toggle="dropdown">更多操作
                      <span class="caret"></span>
                    </button>
                    <ul class="dropdown-menu ssl-more-operations" role="menu">
                        <li>
                            <a href="#">功能</a>
                        </li>
                        <li>
                            <a href="#">另一个功能</a>
                        </li>
                        <li>
                            <a href="#">其他</a>
                        </li>
                        <li class="divider"></li>
                        <li>
                            <a href="#">分离的链接</a>
                        </li>
                    </ul>
                </div>
            </div> -->
        </div>
    </div>
  </div>
  <!-- 订单详情  其他信息 -->
  <div class="row">
      <div class="col-sm-12 col-md-8">
          <div class="card">
              <div class="card-body ssl_detailDiv">
                  <div style="background: #eff2f7;padding: 20px;">
                    <p class="ssl_p">证书信息</p>
                    <form class="form-horizontal" role="form">
                      <div class="row">
                          <?php foreach($Detail['config_options'] as $k => $val): if($val['name_k'] != 'certssl_cert_year'): ?>
                                <div class="form-group ssl_orderDetail col-md-4 col-sm-12">
                                  <label class="control-label"><?php echo $val['name']; ?>：</label>
                                  <div>
                                    <span><?php echo $val['sub_name']; ?></span>
                                  </div>
                                </div>
                              <?php endif; ?>
                          <?php endforeach; ?>
                      </div>
                    </form>
                  </div>
                  <div style="margin-top: 20px;background: #eff2f7;padding: 20px;">
                      <p class="ssl_p">订单信息</p>
                      <form class="form-horizontal" role="form">
                          <div class="form-group ssl_orderDetail">
                            <label class="control-label">订单号：</label>
                            <div>
                              <span><?php echo $Detail['host_data']['ordernum']; ?></span>
                            </div>
                          </div>
                          <div class="row">
                            <div class="form-group ssl_orderDetail col-md-4 col-sm-12">
                              <label class="control-label">订购时间：</label>
                              <div>
                                <span><?php echo date("Y-m-d H:i",!is_numeric($Detail['host_data']['regdate'])? strtotime($Detail['host_data']['regdate']) : $Detail['host_data']['regdate']); ?></span>
                              </div>
                            </div>
                            <div class="form-group ssl_orderDetail col-md-4 col-sm-12">
                              <label class="control-label">签发时间：</label>
                              <div>
                                <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['issus_time']) ? $Detail['host_data']['certssl_orderinfo']['issus_time'] : '-'; ?></span>
                              </div>
                            </div>
                            <div class="form-group ssl_orderDetail col-md-4 col-sm-12">
                              <label class="control-label">到期时间：</label>
                              <div>
                                <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['due_time']) ? $Detail['host_data']['certssl_orderinfo']['due_time'] : '-'; ?></span>
                              </div>
                            </div>
                          </div>
                          <div class="form-group ssl_orderDetail">
                            <label class="control-label">客户备注：</label>
                            <div>
                              <span><?php echo $Detail['host_data']['remark']; ?></span>
                              <i class="bx bx-edit-alt pointer text-primary" data-toggle="modal"
                              data-target="#sslOrderRemark"></i>
                            </div>
                          </div>
                      </form>
                  </div>
                  <?php if($Detail['host_data']['domainstatus'] != 'Pending' && $Detail['host_data']['domainstatus'] != 'Active'): ?>
                  <div style="margin-top: 20px;background: #eff2f7;padding: 20px;">
                    <p class="ssl_p">其他信息</p>
                    <p style="padding-left: 12px;">联系人信息</p>
                    <form class="form-horizontal" role="form">
                      <div class="row">
                        <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                          <label class="control-label">姓名：</label>
                          <div>
                            <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['lastname']) ? $Detail['host_data']['certssl_orderinfo']['lastname'] : '-'; ?><?php echo isset($Detail['host_data']['certssl_orderinfo']['firstname']) ? $Detail['host_data']['certssl_orderinfo']['firstname'] : '-'; ?></span>
                          </div>
                        </div>
                          <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                            <label class="control-label">部门(职位)：</label>
                            <div>
                              <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['position']) ? $Detail['host_data']['certssl_orderinfo']['position'] : '-'; ?></span>
                            </div>
                          </div>
                          <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                            <label class="control-label">电子邮箱：</label>
                            <div>
                              <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['email']) ? $Detail['host_data']['certssl_orderinfo']['email'] : '-'; ?></span>
                            </div>
                          </div>
                          <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                            <label class="control-label">联系电话：</label>
                            <div>
                              <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['telephone']) ? $Detail['host_data']['certssl_orderinfo']['telephone'] : '-'; ?></span>
                            </div>
                          </div>
                      </div>
                    </form>
                      <?php if((isset($Detail['host_data']['cert_pinfo']) && $Detail['host_data']['cert_pinfo']['cert_type'] != 'dv')): ?>
                      <p style="padding-left: 12px;">公司/组织信息</p>
                      <form class="form-horizontal" role="form">
                          <div class="row">
                              <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                                  <label class="control-label">公司名称：</label>
                                  <div>
                                      <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['orgName']) ? $Detail['host_data']['certssl_orderinfo']['orgName'] : '-'; ?></span>
                                  </div>
                              </div>
                              <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                                  <label class="control-label">地址：</label>
                                  <div>
                                      <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['country_name']) ? $Detail['host_data']['certssl_orderinfo']['country_name'] : '-'; ?>、<?php echo isset($Detail['host_data']['certssl_orderinfo']['province']) ? $Detail['host_data']['certssl_orderinfo']['province'] : '-'; ?>、<?php echo isset($Detail['host_data']['certssl_orderinfo']['locality']) ? $Detail['host_data']['certssl_orderinfo']['locality'] : '-'; ?></span>
                                  </div>
                              </div>
                              <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                                  <label class="control-label">邮政编码：</label>
                                  <div>
                                      <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['postalCode']) ? $Detail['host_data']['certssl_orderinfo']['postalCode'] : '-'; ?></span>
                                  </div>
                              </div>
                              <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                                  <label class="control-label">联系电话：</label>
                                  <div>
                                      <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['company_phone']) ? $Detail['host_data']['certssl_orderinfo']['company_phone'] : '-'; ?></span>
                                  </div>
                              </div>
                              <div class="form-group ssl_orderDetail col-md-6 col-sm-12">
                                  <label class="control-label">公司地址：</label>
                                  <div>
                                      <span><?php echo isset($Detail['host_data']['certssl_orderinfo']['address']) ? $Detail['host_data']['certssl_orderinfo']['address'] : '-'; ?></span>
                                  </div>
                              </div>
                          </div>
                        <?php endif; ?>
                  </div>
                  <?php endif; ?>
              </div>
          </div>
      </div>
      <div class="col-sm-12 col-md-4">
        <div class="card">
            <div class="card-body">
                <div>
                    <p class="ssl_p">快速上手</p>
                    <div class="stepDiv">
                        <div class="step
                            <?php if($Detail['host_data']['domainstatus'] == 'Pending'): ?> sslIncompleteStep <?php endif; if($Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Cancelled' || $Detail['host_data']['domainstatus'] == 'Verifiy_Active'): ?> sslCurrentStep <?php endif; if($Detail['host_data']['domainstatus'] != 'Pending' && $Detail['host_data']['domainstatus'] != 'Active' && $Detail['host_data']['domainstatus'] != 'Cancelled' && $Detail['host_data']['domainstatus'] != 'Verifiy_Active'): ?> sslCompletedStep <?php endif; ?>
                            ">
                            <div class="yuan">01</div>
                            <div class="stepContent">
                                <div style="font-size: 16px;font-weight: bold;">申请证书</div>
                                <?php if($Detail['host_data']['domainstatus'] != 'Pending'): if($Detail['host_data']['domainstatus'] == 'Active'): ?>
                                        <div class="stepContent_text">您可以点击签发开始进行证书申请</div>
                                    <?php endif; if($Detail['host_data']['domainstatus'] == 'Cancelled'): ?>
                                        <div class="stepContent_text">您可以点击重新签发开始进行证书申请</div>
                                    <?php endif; if($Detail['host_data']['domainstatus'] == 'Verifiy_Active'): ?>
                                        <div class="stepContent_text">请您尽快添加解析记录，CA扫描认证通过后即可发布证书</div>
                                    <?php endif; if($Detail['host_data']['domainstatus'] == 'Issue_Active' || $Detail['host_data']['domainstatus'] == 'Deleted'): ?>
                                        <div class="stepContent_text">当前证书已于<?php echo isset($Detail['host_data']['certssl_orderinfo']['issus_time']) ? $Detail['host_data']['certssl_orderinfo']['issus_time'] : '-'; ?>签发成功</div>
                                    <?php endif; ?>

                                    <div>
                                        <?php if($Detail['host_data']['domainstatus'] == 'Verifiy_Active'): ?>
                                            <a href="#">
                                                <button type="button" style="min-width:80px;" class="btn btn-ssl-b btn-sm waves-effect waves-light verifyDoamin" data-hostid="<?php echo $Detail['host_data']['id']; ?>">域名验证</button>
                                            </a>
                                        <?php endif; if($Detail['host_data']['domainstatus'] == 'Active'): ?>
                                            <a href="#">
                                              <button type="button" style="width:80px;" class="btn btn-ssl-b btn-sm waves-effect waves-light issueCertBtn">签发</button>
                                            </a>
                                        <?php endif; if($Detail['host_data']['domainstatus'] == 'Cancelled'): ?>
                                            <a href="#">
                                                <button type="button" style="min-width:80px;" class="btn btn-ssl btn-sm waves-effect waves-light issueCertBtn">重新签发</button>
                                            </a>
                                        <?php endif; ?>
                                        <!-- 待验证时才显示中止审核 中止审核的意思就是取消订单 -->
                                        <?php if($Detail['host_data']['domainstatus'] == 'Verifiy_Active'): ?>
                                            <a href="#">
                                                <button type="button" style="min-width:80px;" class="btn btn-ssl btn-sm btn-outline-light" id="cancelIssueOrder">终止审核</button>
                                            </a>
                                        <?php endif; ?>
                                    </div>
                                <?php endif; ?>
                            </div>
                        </div>
                        <div class="step
                                    <?php if($Detail['host_data']['domainstatus'] == 'Overdue_Active' || $Detail['host_data']['domainstatus'] == 'Issue_Active'): ?> sslCurrentStep <?php endif; if($Detail['host_data']['domainstatus'] == 'Pending' || $Detail['host_data']['domainstatus'] == 'Active' || $Detail['host_data']['domainstatus'] == 'Cancelled' || $Detail['host_data']['domainstatus'] == 'Verifiy_Active'): ?> sslIncompleteStep <?php endif; if($Detail['host_data']['domainstatus'] == 'Deleted'): ?> sslCompletedStep <?php endif; ?>
                            ">
                            <div class="yuan">02</div>
                            <div class="stepContent"> 
                                <div style="font-size: 16px;font-weight: bold;">使用证书</div>
                                <?php if($Detail['host_data']['domainstatus'] != 'Pending'): ?>
                                    <div class="stepContent_text">您可以下载证书进行安装</div>
                                    <div>
                                        <a href="#">
                                            <button type="button" style="min-width:80px;" class="btn btn-ssl-b btn-sm waves-effect waves-light" id="downSslCert">下载证书</button>
                                        </a>
                                    </div>
                                <?php endif; ?>
                            </div>
                        </div>
                        <div style="border-color: #fff;" class="step
                            <?php if($Detail['host_data']['domainstatus'] == 'Deleted'): ?> sslCurrentStep <?php endif; if($Detail['host_data']['domainstatus'] != 'Deleted'): ?> sslIncompleteStep <?php endif; ?>
                            ">
                            <div class="yuan">03</div>
                            <div class="stepContent"> 
                                <div style="font-size: 16px;font-weight: bold;">证书过期</div>
                                <?php if($Detail['host_data']['domainstatus'] != 'Pending'): ?>
                                    <div class="stepContent_text">证书到期后重新购买即可！</div>
                                <?php endif; ?>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
      </div>
  </div>
  <!-- 验证 -->
  <div class="modal fade" id="sslVerification" tabindex="-1" role="dialog" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="verificationTitle">域名验证</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body" id="sslVerificationContent">
          <form>
           <div class="form-group d-j-a-center">
              <label class="col-sm-2 control-label">当前域名</label>
              <div class="col-sm-10">
                <div style="display: flex;align-items: center;">
                    <select style="width: 200px;" class="form-control" id="domainCurrent">
                    </select>
                    <div style="margin-left: 20px;">
                      <span style="margin-right:10px;">状态：</span>
                      <span class="domain-info-status">待验证</span>
                    </div>
                </div>
              </div>
            </div>
            <div class="form-group d-j-a-center">
              <label for="verificationMethod" class="col-sm-2 control-label">验证方式</label>
              <div class="col-sm-10">
                <div class="verification-method">
                  <div id="dnsVerification" class="verify_type poitner v-active">dns验证</div>
                  <div id="fileVerification" class="verify_type poitner">文件验证</div>
                </div>
              </div>
            </div>
            <div class="form-group d-j-a-center-operation">
              <div class="col-sm-12">
                <div id="v-operation-text">1. 创建DNS解析记录</div>
                <div style="overflow-x: auto;">
                  <div id="dns-table">
                    <p>选择 DNS 域名授权验证方式，您需要到您的域名解析服务商（如阿里云、新网、腾讯、DNSpod等）提供的系统中进行配置。请根据下列提示，在域名管理处添加一个DNS解析记录，操作完成后，请点击按钮进行检测。</p>
                    <table class="table table-bordered">
                      <tbody>
                        <tr>
                          <td class="td_width" id="record_type">记录类型</td>
                          <td id="copyContent"><span id="record_type_info">CNAME</span> <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="record_typecopy" data-clipboard-action="copy" data-clipboard-target="#record_type"></i></td>
                        </tr>
                        <tr>
                          <td class="td_width" id="host_record">主机记录</td>
                          <td id="copyContent"><span id="host_record_info">CNAME</span><i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="host_recordcopy" data-clipboard-action="copy" data-clipboard-target="#host_record"></i></td>
                        </tr>
                        <tr>
                          <td class="td_width" id="record_value">记录值</td>
                          <td id="copyContent"><span id="record_value_info"></span> <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="record_valuecopy" data-clipboard-action="copy" data-clipboard-target="#record_value"></i></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                  <div id="file-table">
                    <p>如使用文件验证，请先确保自己网站可以正常访问。（ 注：HTTP 使用80端口，HTTPS 使用443端口），根据提示创建txt验证文件，同时上传到网站根目录  ，并且可以通过浏览器访问。请注意windows无法创建以点（ . ）开头的文件夹需要使用命令行 mkdir .well-known</p>
                    <table class="table table-bordered">
                      <tbody>
                        <tr>
                          <td class="td_width">文件位置</td>
                          <td id="copyContent"><span id="file_position"></span><i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="file_positioncopy" data-clipboard-action="copy" data-clipboard-target="#record_value"></i></td>
                        </tr>
                        <tr>
                          <td class="td_width">文件名称</td>
                          <td id="copyContent"><span id="file_name"></span></span> <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="file_namecopy" data-clipboard-action="copy" data-clipboard-target="#record_value"></i></td>
                        </tr>
                        <tr>
                          <td class="td_width" id="record_value">文件内容</td>
                          <td id="copyContent" id="file_verified"><span id="file_verified"></span> <i class="bx bx-copy pointer text-primary ml-1 btn-copy" id="file_verifiedcopy" data-clipboard-action="copy" data-clipboard-target="#record_value"></i></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
            <div class="form-group d-j-a-center">
              <div class="col-sm-12">
                <div id="v-operation-text">2.  检查验证状态</div>
                <div>
                  <span>点击验证按钮检查域名验证状态</span>
                  <a href="#">
                      <button type="button" style="min-width:80px;" class="btn btn-primary btn-sm waves-effect waves-light" id="inspectVerifyStatus">立即检查</button>
                  </a>
                </div>
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button  data-dismiss="modal" id="closesslVerification" type="button" style="width:80px;border-color: #6064FF;color:#6064FF;" class="btn btn-outline-light">关闭</button>
        </div>
      </div>
    </div>
  </div>
  <!-- 签发 -->
  <div class="modal fade" id="sslIssue" tabindex="-1" role="dialog" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="IssueTitle"> 申请签发</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body" id="sslIssueContent">
          <div class="sslIssueDiv">
            <form id="issus_form">
              <div class="row">
                <div class="col-sm-6">
                    <div class="form-group d-j-a-center">
                      <label for="IssueDomain" class="col-sm-3 red_label control-label label-text-right">签发域名</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['used_domainNames']) ? $Detail['host_data']['certssl_orderinfo']['used_domainNames'] : ''; ?>" class="form-control" id="IssueDomain" name="used_domainNames" placeholder="请输入签发域名">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center" id="issueDomainTow" style="">
                      <label for="IssueDomainTow" class="col-sm-3 control-label label-text-right"></label>
                      <div class="col-sm-9">
                        <textarea class="form-control" rows="5" name="domainNames" id="domainVerification" placeholder="不同域名或公网 IP 地址之间请使用换行符隔开，我们支持以下格式：
idcsmart.com
*.idcsmart.com"><?php echo isset($Detail['host_data']['certssl_orderinfo']['domainNames']) ? $Detail['host_data']['certssl_orderinfo']['domainNames'] : ''; ?></textarea>
                        <div style="position: absolute;bottom: 5px;right: 30px;">
                          <span id="currentNum">0</span>/<span id="domainNum">0</span>
                        </div>
                      </div>
                    </div>
                    
                    
                    <div class="form-group d-j-a-center">
                      <label for="issueCsr" class="col-sm-3 red_label control-label label-text-right">CSR</label>
                      <div class="col-sm-9">
                        <div class="issueCsrDiv" style="display: flex;align-items: center;">
                          <div>
                            <input type="radio" name="is_csr" id="issueCsrRadios1" value="0" checked> <span>自动生成</span>
                          </div>
                          <div>
                            <input type="radio" name="is_csr" id="issueCsrRadios2"  value="1"> <span>自行上传</span>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="form-group d-j-a-center" id="RadioTextarea" style="display: none">
                      <label for="issueCsr" class="col-sm-3 control-label"></label>
                      <div class="col-sm-9">
                        <textarea class="form-control" rows="3" name="csr_text" id="algorithmTips"></textarea>
                      </div>
                    </div>
                </div>
                <div class="col-sm-6">
                    <div class="form-group d-j-a-center">
                      <div class="col-sm-8">
                        <!-- 域名验证 -->
                      </div>
                      <div class="col-sm-4">

                      </div>
                    </div>
                </div>
              </div>
              <div style="font-weight: bold;margin-bottom: 5px;padding-left:28px;" id="issueseniorDivTitle">其他信息</div>
              <div class="row" id="issueseniorDiv">
                <div class="col-sm-12">
                  <div style="margin-bottom: 10px;padding-left: 28px;">联系人信息</div>
                  <div class="row verificationDiv">
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueLastname" class="col-sm-3 control-label label-text-right">姓</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['lastname']) ? $Detail['host_data']['certssl_orderinfo']['lastname'] : ''; ?>" class="form-control" name="lastname" id="issueLastname" placeholder="请输入姓">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueName" class="col-sm-3 control-label label-text-right">名</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['firstname']) ? $Detail['host_data']['certssl_orderinfo']['firstname'] : ''; ?>" class="form-control" name="firstname" id="issueName" placeholder="请输入名">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issuePosition" class="col-sm-3 control-label label-text-right">职位</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['position']) ? $Detail['host_data']['certssl_orderinfo']['position'] : ''; ?>" class="form-control" name="position" id="issuePosition" placeholder="请输入职位">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueEmail" class="col-sm-3 control-label label-text-right">邮箱</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['email']) ? $Detail['host_data']['certssl_orderinfo']['email'] : ''; ?>" class="form-control" name="email" id="issueEmail" placeholder="请输入邮箱">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issuePhone" class="col-sm-3 control-label label-text-right">电话</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['telephone']) ? $Detail['host_data']['certssl_orderinfo']['telephone'] : ''; ?>" class="form-control" name="telephone" id="issuePhone" placeholder="请输入电话">
                      </div>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 company-box">
                  <div style="margin-bottom: 10px;padding-left: 28px;">企业信息</div>
                  <div class="row verificationDiv_t">
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCorporatename" class="col-sm-3 redStar control-label label-text-right">公司名称</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['orgName']) ? $Detail['host_data']['certssl_orderinfo']['orgName'] : ''; ?>" class="form-control" name="orgName" id="issueCorporatename" placeholder="请输入公司名称">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCorporatephone" class="col-sm-3 redStar control-label label-text-right">公司电话</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['company_phone']) ? $Detail['host_data']['certssl_orderinfo']['company_phone'] : ''; ?>" class="form-control" id="issueCorporatephone" name="company_phone" placeholder="请输入公司名称">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCreditCode" class="col-sm-3 redStar control-label label-text-right">信用代码</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['creditCode']) ? $Detail['host_data']['certssl_orderinfo']['creditCode'] : ''; ?>" class="form-control" name="creditCode" id="issueCreditCode" placeholder="请输入信用代码">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCountry" class="col-sm-3 redStar control-label label-text-right">国家</label>
                      <div class="col-sm-9">
                        <select class="form-control" id="issueCountry" name="country">
                          <?php foreach($Detail['host_data']['iso_arr'] as $key => $val): ?>
                              <option value="<?php echo $val['iso']; ?>" <?php if((isset($Detail['host_data']['certssl_orderinfo']['country']) && $Detail['host_data']['certssl_orderinfo']['country'] == $val['iso'])): ?> selected <?php endif; if((!isset($Detail['host_data']['certssl_orderinfo']['country']) && $val['iso'] == 'CN')): ?> selected <?php endif; if((isset($Detail['host_data']['certssl_orderinfo']['country']) && empty($Detail['host_data']['certssl_orderinfo']['country']) && $val['iso'] == 'CN')): ?> selected <?php endif; ?>
                              ><?php echo $val['name_zh']; ?>(<?php echo $val['iso']; ?>)</option>
                          <?php endforeach; ?>
                        </select>
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueProvince" class="col-sm-3 redStar control-label label-text-right">省份</label>
                      <div class="col-sm-9">
                          <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['province']) ? $Detail['host_data']['certssl_orderinfo']['province'] : ''; ?>" class="form-control" name="province" id="issueCreditCode" placeholder="请输入省份">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCity" class="col-sm-3 redStar control-label label-text-right">城市</label>
                      <div class="col-sm-9">
                          <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['locality']) ? $Detail['host_data']['certssl_orderinfo']['locality'] : ''; ?>" class="form-control" name="locality" id="issueCreditCode" placeholder="请输入城市">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCompanyaddress" class="col-sm-3 redStar control-label label-text-right">公司地址</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['address']) ? $Detail['host_data']['certssl_orderinfo']['address'] : ''; ?>" class="form-control" id="issueCompanyaddress" name="address" placeholder="请输入公司地址">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issuePostalCode" class="col-sm-3 redStar control-label label-text-right">邮政编码</label>
                      <div class="col-sm-9">
                          <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['postalCode']) ? $Detail['host_data']['certssl_orderinfo']['postalCode'] : ''; ?>" class="form-control" id="issuePostalCode" placeholder="请输入邮政编码" name="postalCode">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueCountryRegistration" class="col-sm-3 redStar control-label label-text-right">注册国家</label>
                      <div class="col-sm-9">
                          <select class="form-control" id="issueCountry" name="joiCountry">
                              <?php foreach($Detail['host_data']['iso_arr'] as $key => $val): ?>
                                  <option value="<?php echo $val['iso']; ?>" <?php if((isset($Detail['host_data']['certssl_orderinfo']['joiCountry']) && $Detail['host_data']['certssl_orderinfo']['joiCountry'] == $val['iso'])): ?> selected <?php endif; if((!isset($Detail['host_data']['certssl_orderinfo']['joiCountry']) && $val['iso'] == 'CN')): ?> selected <?php endif; if((isset($Detail['host_data']['certssl_orderinfo']['joiCountry']) && empty($Detail['host_data']['certssl_orderinfo']['joiCountry']) && $val['iso'] == 'CN')): ?> selected <?php endif; ?>
                                  ><?php echo $val['name_zh']; ?>(<?php echo $val['iso']; ?>)</option>
                              <?php endforeach; ?>
                          </select>
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueRegisteredProvince" class="col-sm-3 redStar control-label label-text-right">注册省份</label>
                      <div class="col-sm-9">
                          <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['joiProvince']) ? $Detail['host_data']['certssl_orderinfo']['joiProvince'] : ''; ?>" class="form-control" name="joiProvince" id="issueCreditCode" placeholder="请输入注册省份">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueRegisteredCity" class="col-sm-3 redStar control-label label-text-right">注册城市</label>
                      <div class="col-sm-9">
                          <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['joiLocality']) ? $Detail['host_data']['certssl_orderinfo']['joiLocality'] : ''; ?>" class="form-control" name="joiLocality" id="issueCreditCode" placeholder="请输入注册城市">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueRegisteredAddress" class="col-sm-3 redStar control-label label-text-right">注册地址</label>
                      <div class="col-sm-9">
                        <input type="text" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['registryAddr']) ? $Detail['host_data']['certssl_orderinfo']['registryAddr'] : ''; ?>" class="form-control" id="issueRegisteredAddress" placeholder="请输入注册地址" name="registryAddr">
                      </div>
                    </div>
                    <div class="form-group d-j-a-center col-sm-12 col-md-6">
                      <label for="issueDateRegistration" class="col-sm-3 redStar control-label label-text-right">注册日期</label>
                      <div class="col-sm-9">
                        <input type="date" value="<?php echo isset($Detail['host_data']['certssl_orderinfo']['dateOfIncorporation']) ? $Detail['host_data']['certssl_orderinfo']['dateOfIncorporation'] : ''; ?>" class="form-control" id="issueDateRegistration" placeholder="请输入注册日期" name="dateOfIncorporation">
                      </div>
                    </div>
                  </div>
                </div>
              </div>
                <input type="hidden" name="id" value="<?php echo $Detail['host_data']['id']; ?>">
                <input type="hidden" name="func" value="issue">
             </form>
          </div>
        </div>
        <div class="modal-footer">
          <button id="issueApply" type="button" class="btn btn-primary" style="width:80px;">提交</button>
          <button id="issueApplyCancel" type="button" class="btn btn-outline-light" style="width:80px;"  data-dismiss="modal">取消</button>
        </div>
      </div>
    </div>
  </div>
  <!-- 修改备注弹窗 -->
  <div class="modal fade" id="sslOrderRemark" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
    aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">修改备注</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <form id="sslOrderRemarkForm" class="needs-validation" novalidate>
            <div class="form-group row">
              <label class="col-sm-3 col-form-label label-text-right">备注</label>
              <div class="col-sm-8">
                <input type="textarea" class="form-control" id="sslOrderRemarkInput" value='<?php echo $Detail['host_data']['remark']; ?>' placeholder="请输入备注" required />
              </div>
            </div>
          </form>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
          <button type="button" class="btn btn-primary mr-2" id="modifyNotesSubmit" onclick="editSubmit(<?php echo $Detail['host_data']['id']; ?>)">确定</button>
        </div>
      </div>
    </div>
  </div>
  <!-- ip弹窗 -->
  <div class="modal fade" id="sslIPdialog" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
    aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">IP地址</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <div id="ipDiv">

              <?php if(($Detail['host_data']['certssl_orderinfo']['domainNames'])): foreach($Detail['host_data']['certssl_orderinfo']['domainNames_arr'] as $k => $v): ?>
                      <p><?php echo $v; ?></p>
                  <?php endforeach; ?>
              <?php endif; ?>

          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
          <button type="button" class="btn btn-primary mr-2" >确定</button>
        </div>
      </div>
    </div>
  </div>
  <script src="/themes/clientarea/default/assets/js/billing.js?v=<?php echo $Ver; ?>"></script>
  <script type="text/javascript">

    $("#domainVerification").on('blur', function () {
      let verificationText=''
      var arrtwo=[]
      var arr=$(this).val().split("\n").map(item =>{
        if(item.replace(/\s+/g,"") != '') arrtwo.push(item)
      })
      var ary=arrtwo
      var nary=ary.sort();
      for(var i=0;i<ary.length;i++){
        if (nary[i]==nary[i+1]){
          // alert("数组重复内容："+nary[i]);
          verificationText+=` ${nary[i]},`
        }
      }
      $('#currentNum').text(arrtwo.length)
      if($('#domainNum').text() == arrtwo.length){
        $(this).removeClass('borderRed')
      }else{
        $(this).addClass('borderRed')
      }
      if(verificationText!='') toastr.error(verificationText.slice(0,verificationText.length-1)+'域名重复');
    })

      let isCompanyInfo = 0
      var hostId = '<?php echo $Detail['host_data']['id']; ?>';
      var domainStatus = '<?php echo $Detail['host_data']['domainstatus']; ?>';
      var timer;

      // 状态为未付款时， 轮询
      if(domainStatus == 'Pending')
      {
          var timer = setInterval(function() {

              $.ajax({
                  url: '/getHostStatus'
                  ,data: {hid:hostId}
                  ,type: 'get'
                  ,dataType: 'json'
                  ,success: function (e) {
                      if(e.status != 200)
                      {
                          clearInterval(timer);
                      }
                      if(e.data == 1)
                      {
                          clearInterval(timer);
                          location.reload();
                      }
                  }
              })
          }, 3000);
      }



      $('.sslIncompleteStep button').attr('disabled', true);
      $('.sslCompletedStep button').attr('disabled', true)
      // 签发提交
      $('#issueApply').click(function(){
          let _that= $(this)
          let text= $(this).text()
          $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2"  style="color:#fff;"></i>' + text)
          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: $('#issus_form').serialize()
              ,type:'post'
              ,dataType: 'json'
              ,success: function(e){
                    _that.html(text)
                    if(e.status != 200)
                    {
                      toastr.error(e.msg);
                      return false;
                    }
                    toastr.success(e.msg);
                    setTimeout(function () {
                      location.reload();
                    }, 2000)
              }
          })

      })
      $('#downSslCert').click(function () {
          let _that= $(this)
          let text= $(this).text()
          $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2" style="color:#fff;"></i>' + text)
          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: {id:hostId, func:'downloadCert'}
              ,type:'post'
              ,success: function(e){
                   _that.html(text)
                  if(e.status != 200)
                  {
                      toastr.error(e.msg);
                      return false;
                  }
                  window.open(e.data.path + hostId );
              }
          })
      })
      // 域名验证 获取验证所需信息
      $('.verifyDoamin').click(function() {
          let text= $(this).text()
          $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2"  style="color:#fff;"></i>' + text)
          getVerifyInfo();
      })
      // 取消审核
      $('#cancelIssueOrder').click(function () {
          let _that= $(this)
          let text= $(this).text()
          $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2"></i>' + text)
          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: {id:hostId, func:'cancelVerify'}
              ,type:'post'
              ,dataType: 'json'
              ,success: function(e){
                  _that.html(text)
                  if(e.status != 200)
                  {
                      toastr.error(e.msg);
                      return false;
                  }
                  toastr.success(e.msg);
                  setTimeout(function () {
                      location.reload();
                  }, 2000)
              }
          })
      })

      /*
       * 签发是 企业信息是否必填，以及 csr 算法的提示信息
       */
      $('.issueCertBtn').click(function () {

          let _that= $(this)
          let text= $(this).text()
          if(_that.hasClass("btn-ssl")){
            $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2"></i>' + text)
          }else{
            $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2" style="color:#fff;"></i>' + text)
          }
          issueBeforeCheckInfo(function (e) {
            _that.html(text)
              if(e.status != 200)
              {
                  toastr.error(e.msg);
                  return false;
              }
              // 企业信息是否必填
              var is_companyInfo = e.data.is_companyInfo
              isCompanyInfo = e.data.is_companyInfo
              for(let item of $('.verificationDiv_t')[0].children){
                 if(is_companyInfo == 0){
                   item.children[0].classList.remove("redStar");
                 }else{
                  item.children[0].classList.add("redStar");
                 }
              }

              if(is_companyInfo == 0)
              {
                  $('.company-box').hide();
              }
              // csr提示信息
              var algorithm_tips = e.data.algorithm_tips
              if(algorithm_tips!=''){
                  $('#algorithmTips')[0].placeholder=algorithm_tips
              }else{
                  $('#algorithmTips')[0].placeholder=''
              }
              $('#sslIssue').modal('show');
              if(e.data.flex_num == 1)
              {
                  $('#issueDomainTow').hide()
              }
              if(e.data.flex_num > 1)
              {
                  $('#issueDomainTow').show()
              }
              $('#domainNum').text(e.data.flex_num - 1)
              var arrtwo=[]
              var arr=$("#domainVerification").val().split("\n").map(item =>{
                if(item.replace(/\s+/g,"") != '') arrtwo.push(item)
              })
              $('#currentNum').text(arrtwo.length)
          })
      });
      // 检查验证状态
      $('#inspectVerifyStatus').click(function () {
          let _that= $(this)
          let text= $(this).text()
          $(this).html('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2"  style="color:#fff;"></i>' + text)
          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: {id:hostId, func:'getVerifiedStatus'}
              ,type:'post'
              ,dataType: 'json'
              ,success: function(e){
                  _that.html(text)
                  if(e.status != 200)
                  {
                      toastr.error(e.msg);
                      return false;
                  }
                  toastr.success(e.msg);
                  setTimeout(function () {
                      location.reload();
                  }, 2000)
              }
          })
      })

      function issueBeforeCheckInfo(fun)
      {   
          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: {id:hostId, func:'issueBeforeCheckInfo'}
              ,type: 'post'
              ,dataType: 'json'
              ,success: function (e) {
                  fun(e)
              }
          })
      }

      /**
       * 获取验证信息
       */
      function getVerifyInfo()
      {
          var domainName = $('#domainCurrent').val();
          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: {id:hostId, func:'getVerifiedInfo', domainName : domainName}
              ,type:'post'
              ,dataType: 'json'
              ,success: function(e){
                $('.verifyDoamin').html('域名验证')
                  if(e.status != 200)
                  {
                      toastr.error(e.msg);
                      return false;
                  }
                  handleVerifyContent(e.data)
              }
          })
      }

      $('#domainCurrent').change(function () {
          getVerifyInfo();
      })

      /**
       * 验证信息回显
       * @param data
       */
      function handleVerifyContent(data)
      {
          var list_html = '';

          for(var o in data.list)
          {
              list_html += '<option value="'+ data.list[o].dnsNames +'"';
              if(data.list[o].dnsNames == data.info.dnsNames)
              {
                  $('.domain-info-status').html(data.list[o].status == 2002 ? '已验证' : '待验证');
                  list_html += 'selected';
              }
              list_html += '>' + data.list[o].dnsNames;
              if(data.list[o].status == 2001)
              {
                  list_html += '-待验证'
              }
              if(data.list[o].status == 2002)
              {
                  list_html += '-已验证'
              }
              list_html += '</option>';
          }

          $('.verify_type').removeClass('v-active');

          if(data.info.dcvMethod == 'dns')
          {
              $('#dnsVerification').addClass('v-active');
              $('#v-operation-text').text('1. 创建DNS解析记录');
              $('#dns-table').css({'display':'block'})
              $('#file-table').css({'display':'none'})
              $('#record_type_info').html(data.info.recordType);
              $('#host_record_info').html(data.info.hostRecord);
              $('#record_value_info').html(data.info.hashValue);

              $('#sslVerification').modal('show');
              $('#domainCurrent').html(list_html);
              return;
          }
          $('#fileVerification').addClass('v-active');
          $('#v-operation-text').text('1. 创建文件解析记录');
          $('#dns-table').css({'display':'none'})
          $('#file-table').css({'display':'block'})
          $('#file_position').html(data.info.filePath);
          $('#file_name').html(data.info.fileName);
          $('#file_verified').html(data.info.hashValue);

          $('#sslVerification').modal('show');
          $('#domainCurrent').html(list_html);
      }
  // 选择验证方式
  $("#dnsVerification").on('click', function () {

      $('.verify_type').removeClass('v-active');
      $(this).addClass('v-active');
      $('#v-operation-text').text('1. 创建DNS解析记录')
      $('#dns-table').css({'display':'block'})
      $('#file-table').css({'display':'none'})

      switchVerifyMode('dns', function (e) {
          if(e.status != 200)
          {
              toastr.error(e.msg);
          }
          getVerifyInfo()
      })
  })
  $("#fileVerification").on('click', function () {

      $('.verify_type').removeClass('v-active');
      $(this).addClass('v-active');
      $('#v-operation-text').text('1. 上传文件进行验证')
      $('#dns-table').css({'display':'none'})
      $('#file-table').css({'display':'block'})

      switchVerifyMode('file', function(e) {
          if(e.status != 200)
          {
              toastr.error(e.msg);
          }
          getVerifyInfo()
      })

  })

      function switchVerifyMode(verifyMethod = 'dns', fun)
      {
          var domainName = $('#domainCurrent').val();

          $.ajax({
              url: '/provision/sslCertFunc'
              ,data: {id:hostId, func:'replaceDcvMethod',verifyMethod: verifyMethod, domainName : domainName}
              ,type:'post'
              ,dataType: 'json'
              ,success: function(e){
                  fun(e)
              }
          })
      }

      // 修改备注 
    function editSubmit(serverid){
      $.ajax({
        type: "POST",
        url: '/host/remark',
        data: {
          id: serverid,
          remark: $('#sslOrderRemarkInput').val()
        },
        dataType: "json",
        success: function (data) {
          toastr.success(data.msg);
          $('#sslOrderRemark').modal('hide')
          location.reload()
        }
      });
    }

    //复制 记录类型
    var clipboard = new ClipboardJS('#record_typecopy', {
      text: function (trigger) {
        console.log($('#record_type_info').text())
        return $('#record_type_info').text()
      },
      container: document.getElementById('sslVerification')   //弹框里的复制 必须获取焦点
    });
    clipboard.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    //复制 主机记录
    var clipboardT = new ClipboardJS('#host_recordcopy', {
      text: function (trigger) {
        return $('#host_record_info').text()
      },
      container: document.getElementById('sslVerification') 
    });
    clipboardT.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    //复制 记录值
    var clipboardS = new ClipboardJS('#record_valuecopy', {
      text: function (trigger) {
        return $('#record_value_info').text()
      },
      container: document.getElementById('sslVerification') 
    });
    clipboardS.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    //复制 文件位置
    var clipboardS = new ClipboardJS('#file_positioncopy', {
      text: function (trigger) {
        return $('#file_position').text()
      },
      container: document.getElementById('sslVerification') 
    });
    clipboardS.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    //复制 文件名
    var clipboardS = new ClipboardJS('#file_namecopy', {
      text: function (trigger) {
        return $('#file_name').text()
      },
      container: document.getElementById('sslVerification') 
    });
    clipboardS.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
    //复制 文件内容
    var clipboardS = new ClipboardJS('#file_verifiedcopy', {
      text: function (trigger) {
        return $('#file_verified').text()
      },
      container: document.getElementById('sslVerification') 
    });
    clipboardS.on('success', function (e) {
      toastr.success('<?php echo $Lang['copy_succeeded']; ?>');
    })
  
  //高级选择
  $("#issueContacts").on('click', function () {
    console.log(123123,$("#issueContacts")[0].checked)
    if($("#issueContacts")[0].checked){
      $('#issueseniorDiv').show()
      $('#issueseniorDivTitle').show()
    }else{
      $('#issueseniorDiv').hide()
      $('#issueseniorDivTitle').hide()
    }
  })
  //csr 操作
  $("#issueCsrRadios2").on('change', function () {
      $('#RadioTextarea').show()

  })
  $("#issueCsrRadios1").on('change', function () {
      $('#RadioTextarea').hide()
  })
  //签发打开弹窗
  $("#IssueBtn").on('click', function () {
    $('#IssueTitle')[0].innerText='申请签发'
  })
  $("#againIssueBtn").on('click', function () {
    $('#IssueTitle')[0].innerText='重新签发'
  })

  //输入姓 验证
  $("#issueLastname").on('blur', function () {
    var dom = $("#issueLastname")[0]
		if (dom.value == '') {
			
			dom.classList.add("is-invalid"); //添加非法状态
			return
		} else {
			dom.classList.remove("is-invalid");
			
		}
  })
  //输入 名验证
  $("#issueName").on('blur', function () {
    var dom = $("#issueName")[0]
		if (dom.value == '') {
			
			dom.classList.add("is-invalid"); //添加非法状态
			return
		} else {
			dom.classList.remove("is-invalid");
			
		}
  })
  //输入职位验证
  $("#issuePosition").on('blur', function () {
    var dom = $("#issuePosition")[0]
		if (dom.value == '') {
			
			dom.classList.add("is-invalid"); //添加非法状态
			return
		} else {
			dom.classList.remove("is-invalid");
			
		}
  })
  //输入邮箱验证
  $("#issueEmail").on('blur', function () {
    var dom = $("#issueEmail")[0]
		if (dom.value == '') {
			
			dom.classList.add("is-invalid"); //添加非法状态
			return
		} else {
			dom.classList.remove("is-invalid");
			
		}
  })
  //输入电话验证
  $("#issuePhone").on('blur', function () {
    var dom = $("#issuePhone")[0]
		if (dom.value == '') {
			
			dom.classList.add("is-invalid"); //添加非法状态
			return
		} else {
			dom.classList.remove("is-invalid");
			
		}
  })
  //输入公司名称验证
  $("#issueCorporatename").on('blur', function () {
    var dom = $("#issueCorporatename")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入公司电话验证
  $("#issueCorporatephone").on('blur', function () {
    var dom = $("#issueCorporatephone")[0]
    if(isCompanyInfo == 1){    
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入信用代码验证
  $("#issueCreditCode").on('blur', function () {
    var dom = $("#issueCreditCode")[0]
    if(isCompanyInfo == 1){    
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入国家验证
  $("#issueCountry").on('blur', function () {
    var dom = $("#issueCountry")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入省份验证
  $("#issueProvince").on('blur', function () {
    var dom = $("#issueProvince")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入城市验证
  $("#issueCity").on('blur', function () {
    var dom = $("#issueCity")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入公司地址验证
  $("#issueCompanyaddress").on('blur', function () {
    var dom = $("#issueCompanyaddress")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入邮政编码验证
  $("#issuePostalCode").on('blur', function () {
    var dom = $("#issuePostalCode")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入注册国家
  $("#issueCountryRegistration").on('blur', function () {
    var dom = $("#issueCountryRegistration")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入注册省份
  $("#issueRegisteredProvince").on('blur', function () {
    var dom = $("#issueRegisteredProvince")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入注册城市
  $("#issueRegisteredCity").on('blur', function () {
    var dom = $("#issueRegisteredCity")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入注册地址
  $("#issueRegisteredAddress").on('blur', function () {
    var dom = $("#issueRegisteredAddress")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  //输入注册日期
  $("#issueDateRegistration").on('blur', function () {
    var dom = $("#issueDateRegistration")[0]
    if(isCompanyInfo == 1){
      if (dom.value == '') {
        
        dom.classList.add("is-invalid"); //添加非法状态
        return
      } else {
        dom.classList.remove("is-invalid");
        
      }
    }
  })
  
  $('#showDomain').on('click', function (){
    console.log(123)
    $(this).css("display","none")
    $('#hideDomain').css("display","flex")
    $('#issueDomainTow').show()
  })
  $('#hideDomain').on('click', function (){
    console.log(123)
    $(this).css("display","none")
    $('#showDomain').css("display","flex")
    $('#issueDomainTow').hide()
  })
  </script>


</body>

</html>
<?php else: ?>
<?php endif; if($ErrorMsg): ?><!-- <div class="alert alert-danger mb-4">
	<div class="alert-body"></div>
</div> -->
<script type="text/javascript">
	$(function () {
		toastr.error('');
	});
</script><?php endif; ?>
<!-- 二次验证 -->
<div class="modal fade" id="secondVerifyModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">二次验证</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<input type="hidden" value="<?php echo $Token; ?>" />
					<input type="hidden" value="closed" name="action" />
					<div class="form-group row mb-4">
						<label class="col-sm-3 col-form-label text-right">验证方式</label>
						<div class="col-sm-8">
							<select class="form-control" class="second_type" name="type" id="secondVerifyType">
								<?php foreach($AllowType as $type): ?>
									<option value="<?php echo $type['name']; ?>"><?php echo $type['name_zh']; ?>：<?php echo $type['account']; ?></option>
								<?php endforeach; ?>
							</select>
						</div>
					</div>
					<div class="form-group row mb-0">
						<label class="col-sm-3 col-form-label text-right">验证码</label>
						<div class="col-sm-8">
							<div class="input-group">
								<input type="text" name="code" id="secondVerifyCode" class="form-control" placeholder="请输入验证码" />
								<div class="input-group-append" id="getCodeBox">
									<button class="btn btn-secondary" id="secondCode" onclick="getSecurityCode()" type="button">获取验证码</button>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary mr-2" id="secondVerifySubmit" onclick="secondVerifySubmitBtn(this)">确定</button>
			</div>
		</div>
	</div>
</div>


<!-- getModalConfirm 确认弹窗 -->
<div class="modal fade" id="confirmModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="confirmBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="confirmSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>
<!-- getModal 自定义body弹窗 -->
<div class="modal fade" id="customModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="customTitle">提示</h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body" id="customBody">

			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" id="customSureBtn">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
	var Userinfo_allow_second_verify = '<?php echo $Userinfo['allow_second_verify']; ?>'
		,Userinfo_user_second_verify = '<?php echo $Userinfo['user']['second_verify']; ?>'
		,Userinfo_second_verify_action_home = <?php echo json_encode($Userinfo['second_verify_action_home']); ?>
		,Login_allow_second_verify = '<?php echo $Login['allow_second_verify']; ?>'
		,Login_second_verify_action_home = <?php echo json_encode($Login['second_verify_action_home']); ?>;
</script>
<script src="/themes/clientarea/default/assets/js/modal.js?v=<?php echo $Ver; ?>"></script>



<div class="modal fade relegation" id="modalUpgradeStepOne" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['upgrade_downgrade']; ?></h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body">

            </div>
        </div>
    </div>
</div>


  <div id="loading-icon" style="display:none">
    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="margin:auto;background:#fff;display:block;" width="200px" height="200px" viewBox="0 0 100 100" preserveAspectRatio="xMidYMid">
      <g>
        <circle cx="73.801" cy="68.263" fill="#93dbe9" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="0s">
          </animateTransform>
        </circle>
        <circle cx="68.263" cy="73.801" fill="#689cc5" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.062s">
          </animateTransform>
        </circle>
        <circle cx="61.481" cy="77.716" fill="#5e6fa3" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.125s">
          </animateTransform>
        </circle>
        <circle cx="53.916" cy="79.743" fill="#3b4368" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.187s">
          </animateTransform>
        </circle>
        <circle cx="46.084" cy="79.743" fill="#191d3a" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.25s">
          </animateTransform>
        </circle>
        <circle cx="38.519" cy="77.716" fill="#d9dbee" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.312s">
          </animateTransform>
        </circle>
        <circle cx="31.737" cy="73.801" fill="#b3b7e2" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.375s">
          </animateTransform>
        </circle>
        <circle cx="26.199" cy="68.263" fill="#93dbe9" r="3">
          <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;360 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s" begin="-0.437s">
          </animateTransform>
        </circle>
        <animateTransform attributeName="transform" type="rotate" calcMode="spline" values="0 50 50;0 50 50" times="0;1" keySplines="0.5 0 0.5 1" repeatCount="indefinite" dur="1.4925373134328357s">
        </animateTransform>
      </g>
    </svg>
</div>

 <div class="modal fade" id="modalUpgradeConfigStepOne" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered " style="max-width:1000px;">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title mt-0" id="myLargeModalLabel"><?php echo $Lang['upgrade_downgrade_options']; ?></h5>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">&times;</span>
                    </button>
                </div>
                <div class="modal-body">
          
                 </div>
              
          </div>
      </div>
  </div>



<!-- 合同蒙层 -->
<?php if($ForceContract): ?>
<div class="contract_mc"></div>
<div class="modal" style="display: block;z-index:999999;width: 100%">
	<div class="modal-dialog modal-dialog-centered">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">合同签订提示</h5>
			</div>
			<div class="modal-body">
                <span>该产品需要签订, 请 <?php if($ForceContract['base'] || $ForceContract['has_contract']): ?><a href="/contract">前往签订</a><?php else: ?><a href="/contracthost?keywords=<?php echo $Detail['host_data']['id']; ?>">前往签订</a><?php endif; if($ForceContract['force']): ?>，逾期<?php echo $ForceContract['suspended']; ?>天未签订会<?php if($ForceContract['suspended_type']=='suspended'): ?>暂停产品服务<?php else: ?>无法访问产品内页<?php endif; ?></span>
                <?php endif; ?>
			</div>
			<div class="modal-footer">
                <?php if($ForceContract['base'] || $ForceContract['has_contract']): ?>
                    <button type="button" class="btn btn-primary" onclick="javascript: location.href='/contract'">确定</button>
                <?php else: ?>
                    <button type="button" class="btn btn-primary" onclick="javascript: location.href='/contracthost?keywords=<?php echo $Detail['host_data']['id']; ?>'">确定</button>
                <?php endif; if($ForceContract['regdated'] == 0 || $ForceContract['force'] == 0): ?>
				<button type="button" class="btn btn-outline-light qd_cancel">取消</button>
                <?php endif; ?>
			</div>
		</div>
	</div>
</div>
<?php endif; ?>
<!-- 服务器内页 修改备注 弹窗 -->
<div class="modal fade" id="modifyRemarkModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
	aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="exampleModalCenterTitle"><?php echo $Lang['modify_remarks']; ?></h5>
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
			</div>
			<div class="modal-body">
				<form>
					<div class="form-group row">
						<label class="col-sm-2 col-form-label text-right"><?php echo $Lang['remarks']; ?></label>
						<div class="col-sm-9">
							<div class="input-group">
								<input id="remarkInp" type="text" value="<?php echo $Detail['host_data']['remark']; ?>" class="form-control api_password"
									placeholder="<?php echo $Lang['please_input']; ?>" />
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-outline-light" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
				<button type="button" class="btn btn-primary" onclick="modifyRemarkSubmit(<?php echo app('request')->get('id'); ?>)"><?php echo $Lang['determine']; ?></button>
			</div>
		</div>
	</div>
</div>
<!-- 重置密码弹窗 --> 
<div class="modal fade" id="moduleResetPass" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
     aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title mt-0"><?php echo $Lang['reset_password']; ?></h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <form id="resetPassword">
                <div class="modal-body">
                    <div class="form-group row mb-4">
                        <label for="horizontal-firstname-input" class="col-md-3 col-form-label d-flex justify-content-end"><?php echo $Lang['password']; ?></label>
                        <div class="col-md-6">
                            <input type="text" class="form-control getPassword" name="password">
                        </div>
                        <div class="col-md-1 fs-18 d-flex align-items-center">
                            <i class="fas fa-dice create_random_pass pointer" onclick="create_random_pass()"></i>
                        </div>
                        <label id="password-error-tip" class="control-label error-tip" for="password"></label>
                    </div>
                    <div class="alert alert-danger" role="alert">
                        <?php echo $Lang['reset_pwd_tip_one']; ?><br>
                        <?php echo $Lang['reset_pwd_tip_two']; ?><br>
                        <?php echo $Lang['reset_pwd_tip_three']; ?>           
                    </div>
                    <div class="form-group row mb-4">
                        <label for="horizontal-firstname-input" class="col-md-3 col-form-label d-flex justify-content-end">强制关机</label>
                        <div class="col-md-6 pt-9">
                            <label>
                                <input type="checkbox" class="mr-1 getForce" id="force" name="force">同意强制关机
                            </label>
                        </div>
                        <label id="force-error-tip" class="control-label error-tip" for="force">请同意强制关机</label>
                    </div>
                </div>
                <input type="hidden" name="func" value="crack_pass">
                <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
            </form>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
                <button type="button" class="btn btn-primary submit" onclick="moduleResetPass($(this))"><?php echo $Lang['determine']; ?></button>
            </div>
        </div>
    </div>
</div>

<!-- 重装系统弹窗 -->
<div class="modal fade" id="moduleReinstall" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
     aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title mt-0"><?php echo $Lang['reinstall_system']; ?></h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <form>
                <div class="modal-body">
                    <div class="row">
                        <div class="col-md-2  d-flex align-items-center justify-content-end">
                            <label class="float-right mb-0"><?php echo $Lang['system']; ?></label>
                        </div>
                        <div class="col-md-5">
                            <div class="form-group mb-0">

                                <select class="form-control configoption_os_group selectpicker" data-style="btn-default" name="os_group" onchange="moduleReinstallOsGroupChange($(this))">
                                  <?php foreach($Detail['cloud_os_group'] as $item): if(strtolower($item['name'])=="windows"): $os_svg = '1'; elseif(strtolower($item['name'])=="centos"): $os_svg = '2'; elseif(strtolower($item['name'])=="ubuntu"): $os_svg = '3'; elseif(strtolower($item['name'])=="debian"): $os_svg = '4'; elseif(strtolower($item['name'])=="esxi"): $os_svg = '5'; elseif(strtolower($item['name'])=="xenserver"): $os_svg = '6'; elseif(strtolower($item['name'])=="freebsd"): $os_svg = '7'; elseif(strtolower($item['name'])=="fedora"): $os_svg = '8'; else: $os_svg = '9'; ?>
                                    <?php endif; ?>
                                    <option  data-content="<img class='mr-1' src='/upload/common/system/<?php echo $os_svg; ?>.svg' height='20'/><?php echo $item['name']; ?>" value="<?php echo $item['id']; ?>"><?php echo $item['name']; ?></option>
                                  <?php endforeach; ?>
                                </select>
                            </div>
                        </div>
						
                        <div class="col-md-5">
                            <div class="form-group">
                                <select class="form-control" name="os" data-os='<?php echo json_encode($Detail['cloud_os']); ?>'>
                                    <?php foreach($Detail['cloud_os'] as $item): ?>
                                       
                                        
                                          <option value="<?php echo $item['id']; ?>" data-group="<?php echo $item['group']; ?>" <?php if($item['group'] != $Detail['cloud_os_group'][0]['id']): ?>style="display:none;"<?php endif; ?> ><?php echo $item['name']; ?></option>
                                        
                                    <?php endforeach; ?>
                                    
                                </select>
                                
                            </div>
                        </div>
                    </div>
                    <?php if($Detail['reinstall_random_port']): ?>
                        <div class="form-group row">
                            <label for="horizontal-firstname-input" class="col-md-2 col-form-label d-flex justify-content-end"><?php echo $Lang['port']; ?></label>
                            <div class="col-md-5">
                                <input type="text" class="form-control" name="port" value="22">
                            </div>
                            <div class="col-md-1 fs-18 d-flex align-items-center">
                                <i class="fas fa-dice module_reinstall_random_port" onclick="module_reinstall_random_port()"></i>
                            </div>
                        </div>
                    <?php endif; if($Detail['reinstall_format_data_disk']): ?>
                    <div class="row">
                        <label class="col-md-2 col-form-label d-flex align-items-center justify-content-end" for="moduleReinstallFormatDataDisk"><?php echo $Lang['format_data_disk']; ?></label>
                        <div class="custom-control custom-switch custom-switch-md mb-4 ml-2" dir="ltr">
                            <input type="checkbox" class="custom-control-input" id="moduleReinstallFormatDataDisk" name="format_data_disk" value="1">
                            <label class="custom-control-label" for="moduleReinstallFormatDataDisk"></label>
                        </div>
                    </div>
                    <?php endif; ?>
                    <div class="row">
                        <div class="col-md-3 offset-md-2 d-flex d-flex-cl justify-content-end">
                            <div class="custom-control custom-checkbox">
                                <input type="checkbox" class="custom-control-input" id="moduleReinstallConfirm" value="1" onchange="reinstallConfirm($(this))">
                                <label class="custom-control-label" for="moduleReinstallConfirm"><?php echo $Lang['finished_backup']; ?></label>
                            </div>
                            <div class="reinstallAgreeCheckbox" id="reinstallAgreeCheckbox"></div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-9">
                            <div id="moduleReinstallMsg"></div>
                        </div>
                    </div>
                </div>
                <input type="hidden" name="func" value="reinstall">
                <input type="hidden" name="id" value="<?php echo app('request')->get('id'); ?>">
            </form>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal"><?php echo $Lang['cancel']; ?></button>
                <button type="button" class="btn btn-primary submit"  onclick="moduleReinstall($(this))"><?php echo $Lang['determine']; ?></button>
            </div>
        </div>
    </div>
</div>

<!-- 升降级商品 -->
<div id="upgradeProductDiv"></div>
<div id="upgradeConfigDiv"></div>
<div id="orderFlowDiv"></div>
<div id="renewDiv"></div>

<script>
    // 重置密码密码校验规则
    var passwordRules =  <?php echo json_encode($Detail['host_data']['password_rule']['rule']); ?>;
    var setting_web_url = '';
    //console.log('passwordRules',passwordRules)
    $(document).on('blur', '.getPassword', function(){
		veriPassword()
    })
    function veriPassword(){
      
        let result = checkingPwd1($(".getPassword").val(), passwordRules.num, passwordRules.upper, passwordRules.lower, passwordRules.special)
      
		if(result.flag) {
			$('#password-error-tip').css('display','none');
			$('.getPassword').removeClass("is-invalid");
		}else{
			$("#password-error-tip").html(result.msg);
			$(".getPassword").addClass("is-invalid");
			$('#password-error-tip').css('display','block');
		}
    }
    $(function(){
        $("#resetPassword").on('click',".create_random_pass",function(e){
            veriPassword()
        })
        $('.qd_cancel').on('click',function(){
            $('.contract_mc,.modal').hide()
        })
    })

    $(".getForce").change(function () {
        if($(this).is(':checked')){
            $('#force-error-tip').css('display', 'none');
        }else{
            $("#force-error-tip").html('请勾选同意强制关机');
            $('#force-error-tip').css('display', 'block');
        }
    });
</script>
<script src="/themes/clientarea/default/assets/js/servicedetail.js?v=<?php echo $Ver; ?>"></script><?php if($TplName != 'login' && $TplName != 'register' && $TplName != 'pwreset' && $TplName != 'bind' && $TplName != 'loginaccesstoken'): ?>
</div>
</div>
</div>

<footer class="footer">
	<div class="container-fluid">
		<div class="row">
			<div class="col-sm-6">
				&copy; <?php echo $Setting['company_name']; ?>.
			</div>
			<div class="col-sm-6">
				<div class="text-sm-right d-none d-sm-block">
					
				</div>
			</div>
		</div>
	</div>
</footer>
<?php endif; ?>
<script src="/themes/clientarea/default/assets/js/app.js?v=<?php echo $Ver; ?>"></script>
<?php $hooks=hook('client_area_footer_output'); if($hooks): foreach($hooks as $item): ?>
		<?php echo $item; ?>
	<?php endforeach; ?>
<?php endif; ?>
</body>

</html>