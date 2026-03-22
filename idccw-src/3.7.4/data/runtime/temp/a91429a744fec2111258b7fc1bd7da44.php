<?php /*a:5:{s:98:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/cart/default/configureproduct.tpl";i:1737452904;s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/head.tpl";i:1737452904;s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/menu.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/pageheader.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/breadcrumb.tpl";i:1737452904;}*/ ?>

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
	.was-validated .custom-control-input:valid~.custom-control-label::before {
		background-color: #fff;
	}

	.bootstrap-select.is-valid .dropdown-toggle,
	.was-validated .bootstrap-select select:valid+.dropdown-toggle {
		border-color: #ced4da !important;
	}

	.form-control.is-valid,
	.was-validated .form-control:valid {
		display: block;
		width: 100%;
		padding: 0.47rem 0.75rem;
		font-size: 0.8125rem;
		font-weight: 400;
		line-height: 1.5;
		color: #495057;
		background-color: #fff;
		background-clip: padding-box;
		border: 1px solid #ced4da;
		border-radius: 0.25rem;
		background-image: url(data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiP…45NSwyLjI2LDYsMi4xMyw2LDJWMS41QzYsMS4yMiw1Ljc4LDEsNS41LDF6Ii8+Cjwvc3ZnPgo=);
		background-position-y: 50%;
		background-position-x: calc(100% - 11px);
		background-repeat: no-repeat;
		padding-right: 24px;
	}

	.cycle-select-discount {
		position: absolute;
		top: -10px;
		right: 0;
		display: inline-block;
		padding: 0px 5px;
		height: 15px;
		background: #F70302;
		border-radius: 5px 5px 5px 0px;
		color: #ffffff;
		font-size: 12px;
		text-align: center;
		line-height: 15px;
	}

	.error-tip {
		color: #f46a6a;
		margin: 0;
		padding: 0;
		line-height: 36px;
		display: none;
	}

	input::-webkit-input-placeholder {
		color: #cacdd4 !important;
	}

	input:-moz-placeholder {
		color: #cacdd4 !important;
	}

	input::-moz-placeholder {
		color: #cacdd4 !important;
	}

	input:-ms-input-placeholder {
		color: #cacdd4 !important;
	}

	.btn-custom {
		min-width: 72px;
		height: 28px;
		border: #F0F0F0 1px solid;
		border-radius: 4px;
		line-height: 11px;
	}

	.btn-custom-group .active {
		background-color: #2948DF;
		color: #fff;
	}

	/* 重写进度条 */
	input[type='range'] {
		background: #F1F3F8;
		outline: none;
		-webkit-appearance: none;
		/*清除系统默认样式*/
		height: 4px;
		/*横条的高度*/
		border-radius: 3px;
		background: rgb(41, 72, 223) !important;
	}

	input[type="range"]::-webkit-slider-thumb {
		-webkit-appearance: none;
		width: 10px;
		height: 23px;
		background-color: #fff;
		/*box-shadow: 0 0 2px rgba(0, 0, 0, 0.3),
		0 3px 5px rgba(0, 0, 0, 0.2);*/
		cursor: pointer;
		border: 4px solid #2948DF;
		border-top-width: 5px;
		border-bottom-width: 5px;
		border-radius: 2px;
	}

	input[type="range"]::-moz-range-thumb {
		-webkit-appearance: none;
		width: 2px;
		height: 15px;
		background-color: #fff;
		/*box-shadow: 0 0 2px rgba(0, 0, 0, 0.3),
		 0 3px 5px rgba(0, 0, 0, 0.2);*/
		cursor: pointer;
		border: 4px solid #2948DF;
		border-top-width: 5px;
		border-bottom-width: 5px;
		border-radius: 2px;
	}

	.range_none {
		position: absolute;
		height: 3px;
		display: block;
		background: #DEDEDE;
		cursor: not-allowed
	}

	.form-check {
		margin-left: 20px;
	}

	.justify-content-start {
		flex-wrap: wrap;
	}

	@media (max-width: 992px) {
		.footer {
			z-index: 1;
			display: none !important;
		}
	}
</style>
<?php if(isset($Get['i'])): ?>
<form id="addCartForm" method="post" class="needs-validation configoption_form" novalidate
	action="?action=configureproduct&pid=<?php echo $CartConfig['product']['id']; ?>&pos[]=<?php echo app('request')->get('i'); ?>"><input type="hidden" name="token" value="<?php echo $Token; ?>">
	<input type="hidden" name="i" value="<?php echo app('request')->get('i'); ?>" />
	<?php else: ?>
	<form id="addCartForm" method="post" class="needs-validation configoption_form" novalidate
		action="?action=configureproduct&pid=<?php echo $CartConfig['product']['id']; if($Get['site']): ?>&site=<?php echo $Get['site']; ?><?php endif; ?>"><input type="hidden" name="token" value="<?php echo $Token; ?>">
		<?php endif; ?>
		<div class="row">
			<div class="col-xl-8">
				<div class="card border-bottom mb-0">
					<div class="card-body">
						<?php if($ErrorMsg): ?>
						<div class="alert alert-danger">
							<a href="#" class="close" data-dismiss="alert">
								&times;
							</a>
							<strong><?php echo $ErrorMsg; ?></strong>
						</div>
						<?php endif; ?>
						<div class="d-flex justify-content-between align-items-center">
							<h4 class="card-title "><?php echo $CartConfig['product']['name']; ?></h4>
						</div>
					</div>
				</div>

				<div class="card">
					<div class="card-body">



						<input type="hidden" name="pid" value="<?php echo $CartConfig['product']['id']; ?>" />
						<input type="hidden" name="currencyid" value="<?php echo $CartConfig['dafault_currencyid']; ?>" />
						<input name="qty" type="hidden" value="1" />
						<?php if((isset($addParam['promocode']))): ?>
						<input type="hidden" name="promocode" value="<?php echo $addParam['promocode']; ?>" />
						<?php endif; if((isset($addParam['aff']))): ?>
						<input type="hidden" name="aff" value="<?php echo $addParam['aff']; ?>" />
						<?php endif; if((isset($addParam['sale']))): ?>
						<input type="hidden" name="sale" value="<?php echo $addParam['sale']; ?>" />
						<?php endif; ?>
						<!--产品配置项-->
						<?php foreach($CartConfig['option'] as $option): if($option['option_type']==1): ?>
						<div class="form-group row configureproduct">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
								<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
									<i class="bx bxs-help-circle pointer text-primary"></i>
								</span>
								<?php endif; ?>
							</label>
							<div class="col-md-3">
								<select id="config<?php echo $option['id']; ?>" name="configoption[<?php echo $option['id']; ?>]" class="form-control selectpicker"
									data-style="btn-default">
									<?php foreach($option['sub'] as $sub): ?>
									<option id="sub<?php echo $sub['id']; ?>" <?php if($CartConfig['config_options'][$option['id']]==$sub['id']): ?> selected="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub['id'] && !$CartConfig['config_options']): ?> selected="" <?php endif; ?>
										value="<?php echo $sub['id']; ?>"><?php echo $sub['option_name']; ?></option>
									<?php endforeach; ?>
								</select>
							</div>
						</div>

						<?php elseif($option['option_type']==20): ?>

						<div class="form-group row configureproduct lingAge-<?php echo $option['id']; ?>">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
								<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
									<i class="bx bxs-help-circle pointer text-primary"></i>
								</span>
								<?php endif; ?>
							</label>
							<div class="col-md-10">
								<div class="btn-group btn-group-toggle mt-2 mt-xl-0" data-toggle="buttons">
									<?php foreach($option['sub'] as $sub_key=>$sub_val): ?>
									<label class="btn btn-primary btn-sm <?php if($CartConfig['config_options'][$option['id']]==$sub_val['id']): ?>active <?php elseif((($option['checkSubId']==$sub_val['id']) && !$CartConfig['config_options'])): ?> active <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?> active
									<?php endif; ?>">
											<input id="config<?php echo $option['id']; ?>_<?php echo $sub_val['id']; ?>" type="radio" data-optionid="<?php echo $option['id']; ?>"
												data-subid="<?php echo $sub_val['id']; ?>" name="configoption[<?php echo $option['id']; ?>]" value="<?php echo $sub_val['id']; ?>"
												class="form-check-input" <?php if($CartConfig['config_options'][$option['id']]==$sub_val['id']): ?>checked="" <?php elseif((($option['checkSubId']==$sub_val['id']) && !$CartConfig['config_options'])): ?> checked="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?> checked=""
												<?php endif; ?>> <?php echo $sub_val['option_name']; ?>
										</label>
									<?php endforeach; ?>
								</div>
							</div>
						</div>
						<div class="lingAge-<?php echo $option['id']; ?>-son">
							<?php if((isset($option['son']) && $option['son'])): foreach($option['son'] as $son_k1 => $son_v1): ?>
							<div class="form-group row configureproduct">
								<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $son_v1['option_name']; if($son_v1['notes']): ?>
									<span data-toggle="tooltip" data-placement="right" title="<?php echo $son_v1['notes']; ?>">
										<i class="bx bxs-help-circle pointer text-primary"></i>
									</span>
									<?php endif; ?>
								</label>
								<div class="col-md-10">
									<div class="btn-group btn-group-toggle mt-2 mt-xl-0" data-toggle="buttons">
										<?php foreach($son_v1['sub'] as $sub_key=>$sub_val): ?>
											<label class="btn btn-primary btn-sm <?php if($CartConfig['config_options'][$son_v1['id']]==$sub_val['id']): ?>active
										<?php elseif((($son_v1['checkSubId']==$sub_val['id']) && !$CartConfig['config_options'])): ?> active <?php elseif($getUrlConfig['config_options'][$son_v1['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?> active
										<?php endif; ?>">
												<input id="config<?php echo $son_v1['id']; ?>_<?php echo $sub_val['id']; ?>" type="radio" data-optionid="<?php echo $option['id']; ?>"
													data-subid="<?php echo $sub_val['id']; ?>" name="configoption[<?php echo $son_v1['id']; ?>]" value="<?php echo $sub_val['id']; ?>"
													class="form-check-input" <?php if($CartConfig['config_options'][$son_v1['id']]==$sub_val['id']): ?>checked=""
													<?php elseif((($son_v1['checkSubId']==$sub_val['id']) && !$CartConfig['config_options'])): ?> checked="" <?php elseif($getUrlConfig['config_options'][$son_v1['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?> checked=""
													<?php endif; ?>> <?php echo $sub_val['option_name']; ?>
											</label>
										<?php endforeach; ?>
									</div>
								</div>
							</div>
							<?php endforeach; ?>
							<?php endif; ?>
						</div>

						<?php elseif($option['option_type']==2): ?>
						<div class="form-group row configureproduct">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
								<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
									<i class="bx bxs-help-circle pointer text-primary"></i>
								</span>
								<?php endif; ?>
							</label>
							<div class="col-md-10 d-flex justify-content-start" style="padding-top: calc(0.47rem + 1px);">
								<?php foreach($option['sub'] as $sub_key=>$sub_val): ?>
								<div class="form-check mr-4">
									<input id="config<?php echo $option['id']; ?>_<?php echo $sub_val['id']; ?>" type="radio" name="configoption[<?php echo $option['id']; ?>]"
										value="<?php echo $sub_val['id']; ?>" class="form-check-input" <?php if($CartConfig['config_options'][$option['id']]==$sub_val['id']): ?>checked="" <?php elseif($sub_key==0 &&
										!$CartConfig['config_options']): ?> checked="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?> checked=""
										<?php endif; ?>> <label class="form-check-label"
										for="config<?php echo $option['id']; ?>_<?php echo $sub_val['id']; ?>"><?php echo $sub_val['option_name']; ?></label>
								</div>
								<?php endforeach; ?>
							</div>
						</div>
						<?php elseif($option['option_type']==3): ?>
						<div class="form-group row configureproduct">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
								<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
									<i class="bx bxs-help-circle pointer text-primary"></i>
								</span>
								<?php endif; ?>
							</label>
							<div class="col-md-10">
								<?php foreach($option['sub'] as $sub_key=>$sub_val): ?>
								<div class="custom-control custom-checkbox mb-3">
									<input id="config<?php echo $option['id']; ?>_<?php echo $sub_val['id']; ?>" type="checkbox" name="configoption[<?php echo $option['id']; ?>]"
										class="custom-control-input" <?php if($CartConfig['config_options'][$option['id']]==$sub_val['id']): ?> checked=""
										<?php elseif($getUrlConfig['config_options'][$option['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?>
										checked="" <?php endif; ?> value="<?php echo $sub_val['id']; ?>">
									<label class="custom-control-label"
										for="config<?php echo $option['id']; ?>_<?php echo $sub_val['id']; ?>"><?php echo $sub_val['option_name']; ?></label>
								</div>
								<?php endforeach; ?>
							</div>
						</div>
						<?php elseif($option['option_type']==4 || $option['option_type']==7 || $option['option_type']==9 ||
						$option['option_type']==11 || $option['option_type']==14 || $option['option_type']==15 || $option['option_type']==16 ||
						$option['option_type']==17 || $option['option_type']==18 || $option['option_type']==19): ?>
						<div class="form-group row configureproduct">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
								<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
									<i class="bx bxs-help-circle pointer text-primary"></i>
								</span>
								<?php endif; ?>
							</label>
							<div class="col-md-10 d-flex align-items-center" style="padding:0px;display: flex;position: relative;">
							<input type="range" min="<?php echo $option['qty_minimum']; ?>" max="<?php echo $option['qty_maximum']; ?>"
								qty_stage="<?php echo $option['qty_stage']==0 ? 1  : $option['qty_stage']; ?>" <?php if($CartConfig['config_options'][$option['id']]): ?>value="<?php echo $CartConfig['config_options'][$option['id']]; ?>" <?php elseif($getUrlConfig['config_options'][$option['id']] && !$CartConfig['config_options']): ?>
								value="<?php echo $getUrlConfig['config_options'][$option['id']]; ?>" <?php else: ?>value="<?php echo $option['qty_minimum']; ?>" <?php endif; ?>
								data-sub='<?php echo json_encode($option['sub']); ?>' class="form-control-range configoption_range float-left mr-2"
								style="width: 100%;" onclick="rangeChange(this)">
								<!--<input type="text" min="<?php echo $option['qty_minimum']; ?>" max="<?php echo $option['qty_maximum']; ?>" <?php if($CartConfig['config_options'][$option['id']]): ?>value="<?php echo $CartConfig['config_options'][$option['id']]; ?>" <?php elseif($getUrlConfig['config_options'][$option['id']] && !$CartConfig['config_options']): ?> value="<?php echo $getUrlConfig['config_options'][$option['id']]; ?>" <?php else: ?>value="<?php echo $option['qty_minimum']; ?>"<?php endif; ?> data-sub='<?php echo json_encode($option['sub']); ?>'  class="configoption_range float-left"> -->
								<input id="config<?php echo $option['id']; ?>" data-type="number"
									class="mr-2 col-md-1 form-control form-control-sm configoption_range_val float-left" style="width:20%"
									name="configoption[<?php echo $option['id']; ?>]" onblur="numberKeyup(this)" type="text" min="<?php echo $option['qty_minimum']; ?>"
									max="<?php echo $option['qty_maximum']; ?>" qty_stage="<?php echo $option['qty_stage']==0 ? 1  : $option['qty_stage']; ?>" <?php if($CartConfig['config_options'][$option['id']]): ?>value="<?php echo $CartConfig['config_options'][$option['id']]; ?>" <?php elseif($getUrlConfig['config_options'][$option['id']] && !$CartConfig['config_options']): ?>
									value="<?php echo $getUrlConfig['config_options'][$option['id']]; ?>" <?php else: ?>value="<?php echo $option['qty_minimum']; ?>" <?php endif; ?>>
								<span><?php echo $option['unit']; ?></span>
									<!--<?php if($option['option_type'] == '4' || $option['option_type'] == '15'): ?>
								<span>个</span>
								<?php elseif($option['option_type'] == '7' || $option['option_type'] == '16'): ?>
								<span>核</span>
								<?php elseif($option['option_type'] == '9' || $option['option_type'] == '17'): ?>
								<span>GB</span>
								<?php elseif($option['option_type'] == '11' || $option['option_type'] == '18'): ?>
								<span>Mbps</span>
								<?php elseif($option['option_type'] == '14' || $option['option_type'] == '19'): ?>
								<span>GB</span>
								<?php endif; ?>-->
								</div>
							</div>
							<?php elseif($option['option_type']==6 || $option['option_type']==8 || $option['option_type']==10 ||
							$option['option_type']==13): ?>
							<div class="form-group row configureproduct">
								<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
									<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
										<i class="bx bxs-help-circle pointer text-primary"></i>
									</span>
									<?php endif; ?>
								</label>
								<div class="col-md-10">
									<div class="btn-group btn-group-toggle mt-2 mt-xl-0" data-toggle="buttons">
										<?php foreach($option['sub'] as $sub_key=>$sub_val): ?>
										<label class="btn btn-primary btn-sm ">
											<input id="config<?php echo $option['id']; ?>_<?php echo $sub_val['id']; ?>" type="radio" <?php if($CartConfig['config_options'][$option['id']]==$sub_val['id']): ?>checked="" <?php elseif($sub_key==0 &&
												!$CartConfig['config_options']): ?> checked="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub_val['id'] && !$CartConfig['config_options']): ?> checked=""
												<?php endif; ?> name="configoption[<?php echo $option['id']; ?>]" value="<?php echo $sub_val['id']; ?>"> <?php echo $sub_val['option_name']; ?>
										</label>
										<?php endforeach; ?>
									</div>
								</div>
							</div>
							<!--操作系统-->
							<?php elseif($option['option_type']==5): ?>
							<div class="form-group row configureproduct">
								<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
									<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
										<i class="bx bxs-help-circle pointer text-primary"></i>
									</span>
									<?php endif; ?>
								</label>
								<?php if($option['sub']['os']): ?>
								<div class="col-md-3">
									<select id="config<?php echo $option['id']; ?>" name="configoption[<?php echo $option['id']; ?>]"
										class="form-control selectpicker configoption_os selectpicker_refresh" data-style="btn-default">
										<?php foreach($option['sub']['os']['child'] as $sub): ?>
										<option id="sub<?php echo $sub['id']; ?>" value="<?php echo $sub['id']; ?>" <?php if($CartConfig['config_options'][$option['id']]==$sub['id']): ?>selected="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub['id'] && !$CartConfig['config_options']): ?> selected=""
											<?php endif; ?>><?php echo $sub['version']; ?> </option> <?php endforeach; ?> </select> </div> <?php else: ?> <div class="col-md-2">
											<select class="form-control configoption_os_group selectpicker" data-style="btn-default"
												onchange="osGroupChange(this)">
												<?php foreach($option['sub'] as $sub_key=>$sub_val): $os_selected = ''; foreach($sub_val['child'] as $child_key=>$child_val): if($child_val['id']==$CartConfig['config_options'][$option['id']]): $os_selected = $sub_key; ?>
												<?php endif; ?>
												<?php endforeach; if(strtolower($sub_key)=="windows"): $os_svg = '1'; elseif(strtolower($sub_key)=="centos"): $os_svg = '2'; elseif(strtolower($sub_key)=="ubuntu"): $os_svg = '3'; elseif(strtolower($sub_key)=="debian"): $os_svg = '4'; elseif(strtolower($sub_key)=="esxi"): $os_svg = '5'; elseif(strtolower($sub_key)=="xenserver"): $os_svg = '6'; elseif(strtolower($sub_key)=="freebsd"): $os_svg = '7'; elseif(strtolower($sub_key)=="fedora"): $os_svg = '8'; else: $os_svg = '9'; ?>
												<?php endif; ?>
												<option
													data-content="<img class='mr-1' src='/upload/common/system/<?php echo $os_svg; ?>.svg' height='20'/><?php echo $sub_key; ?>"
													<?php if($os_selected): ?>selected="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$child_val['id']
													&& !$CartConfig['config_options']): ?> selected="" <?php endif; ?> data-os='<?php echo json_encode($sub_val['child']); ?>'
													value="<?php echo $sub_key; ?>">
													<?php echo $sub_key; ?></option>
												<?php endforeach; ?>
											</select>
								</div>
								<div class="col-md-3">
									<select id="config<?php echo $option['id']; ?>" name="configoption[<?php echo $option['id']; ?>]"
										class="form-control selectpicker_refresh configoption_os" <?php if($CartConfig['config_options'][$option['id']]): ?> data-os-selected="<?php echo $CartConfig['config_options'][$option['id']]; ?>"
										<?php elseif($getUrlConfig['config_options'][$option['id']] && !$CartConfig['config_options']): ?>
										data-os-selected="<?php echo $getUrlConfig['config_options'][$option['id']]; ?>" <?php endif; ?> data-style="btn-default">

									</select>
								</div>
								<?php endif; ?>

							</div>

							<!--数据中心-->
							<?php elseif($option['option_type']==12): ?>
							<div class="form-group row configureproduct">
								<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $option['option_name']; if($option['notes']): ?>
									<span data-toggle="tooltip" data-placement="right" title="<?php echo $option['notes']; ?>">
										<i class="bx bxs-help-circle pointer text-primary"></i>
									</span>
									<?php endif; ?>

								</label>
								<div class="col-md-10">

									<?php if($option['sub']['0']['area']['0']['area_zh']): ?>

									<select id="config<?php echo $option['id']; ?>" name="configoption[<?php echo $option['id']; ?>]"
										class="form-control col-md-3 selectpicker" data-style="btn-default">
										<?php foreach($option['sub'] as $sub_key=>$sub_val): foreach($sub_val['area'] as $area): ?>
										<option id="sub<?php echo $area['id']; ?>"
											data-content="<img class='mr-1' src='/upload/common/country/<?php echo $sub_val['country_code']; ?>.png' height='20'/><?php echo $sub_val['option_name']; ?><?php echo $area['area_zh']; ?>"
											value="<?php echo $area['id']; ?>" <?php if($CartConfig['config_options'][$option['id']]==$area['id']): ?>selected="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$area['id'] && !$CartConfig['config_options']): ?> selected=""
											<?php endif; ?>> <?php echo $sub_val['option_name']; ?><?php echo $area['area_zh']; ?></option> <?php endforeach; ?> <?php endforeach; ?> </select> <?php else: ?> <div
											class=" btn-group-toggle mt-2 mt-xl-0" data-toggle="buttons">
											<?php foreach($option['sub'] as $sub_key=>$sub_val): ?>
											<label class="btn btn-primary   mr-2">
												<input id="config<?php echo $option['id']; ?>_<?php echo $sub_val['area']['0']['id']; ?>" type="radio" <?php if($CartConfig['config_options'][$option['id']]==$sub_val['area']['0']['id']): ?>checked="" <?php elseif($sub_key==0 &&
													!$CartConfig['config_options']): ?> checked="" <?php elseif($getUrlConfig['config_options'][$option['id']]==$sub_val['area']['0']['id'] && !$CartConfig['config_options']): ?>
													checked="" <?php endif; ?> name="configoption[<?php echo $option['id']; ?>]" value="<?php echo $sub_val['area']['0']['id']; ?>">
												<?php if($sub_val['country_code']): ?>
												<img class='mr-1' src='/upload/common/country/<?php echo $sub_val['country_code']; ?>.png' height='20' />
												<?php endif; ?>
												<?php echo !empty($sub_val['option_name']) ? $sub_val['option_name'] : $Lang['defaults']; ?>
											</label>
											<?php endforeach; ?>
								</div>
								<?php endif; ?>

							</div>
						</div>
						<?php endif; ?>



						<?php endforeach; ?>
						<!--自定义字段-->
						<?php foreach($CartConfig['custom_fields'] as $custom_fields): if($custom_fields['fieldtype']=="dropdown"): ?>
						<div class="form-group row">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $custom_fields['fieldname']; ?></label>
							<div class="col-md-3">
								<select id="customfield[<?php echo $custom_fields['id']; ?>]" name="customfield[<?php echo $custom_fields['id']; ?>]"
									class="form-control selectpicker" data-style="btn-default" <?php if($custom_fields['required']): ?>required<?php endif; ?>>
									<?php foreach($custom_fields['dropdown_option'] as $dropdown_option): ?> <option value="<?php echo $dropdown_option; ?>" <?php if($CartConfig['custom_fields_value'][$custom_fields['id']]==$dropdown_option): ?>selected="" <?php endif; ?>>
									<?php echo $dropdown_option; ?></option> <?php endforeach; ?> </select> </div> </div> <?php elseif($custom_fields['fieldtype']=="tickbox"): ?> <div class="form-group row">
									<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $custom_fields['fieldname']; ?></label>
									<div class="col-md-10">
										<div class="custom-control custom-checkbox mb-3">
											<input name="customfield[<?php echo $custom_fields['id']; ?>]" <?php if($CartConfig['custom_fields_value'][$custom_fields['id']]==1): ?>checked="" <?php endif; ?> type="checkbox"
												class="custom-control-input" id="customfields<?php echo $custom_fields['id']; ?>">
											<label class="custom-control-label"
												for="customfields<?php echo $custom_fields['id']; ?>"><?php echo $custom_fields['description']; ?></label>
										</div>
									</div>
							</div>
							<?php elseif($custom_fields['fieldtype']=="textarea"): ?>
							<div class="form-group row">
								<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $custom_fields['fieldname']; ?></label>
								<div class="col-md-3">
									<textarea id="customfield[<?php echo $custom_fields['id']; ?>]" name="customfield[<?php echo $custom_fields['id']; ?>]"
										class="form-control getTextareaVal"
										data-object="<?php echo $custom_fields['required']; ?>||<?php echo $custom_fields['regexpr']; ?>" rows="2"
										placeholder="<?php echo $custom_fields['description']; ?>" <?php if($custom_fields['required']): ?>required<?php endif; ?>><?php echo $CartConfig['custom_fields_value'][$custom_fields['id']]; ?>
										</textarea> </div> </div> <?php else: ?> <div class="form-group row">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $custom_fields['fieldname']; ?></label>
							<div class="col-md-3">
								<input class="form-control getVal" data-object="<?php echo $custom_fields['required']; ?>||<?php echo $custom_fields['regexpr']; ?>" id="customfield[<?php echo $custom_fields['id']; ?>]" name="customfield[<?php echo $custom_fields['id']; ?>]"
									value="<?php echo $CartConfig['custom_fields_value'][$custom_fields['id']]; ?>" type="<?php echo $custom_fields['fieldtype']; ?>"
									placeholder="<?php echo $custom_fields['description']; ?>" <?php if($custom_fields['required']): ?>required<?php endif; ?> />
							</div>
						</div>
						<?php endif; ?>
						<?php endforeach; ?>
						<div class="form-group row">
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $Lang['cycles']; ?></label>
							<div class="col-md-10">
								<div class="btn-group btn-group-toggle cycle mt-2 mt-xl-0 btn-custom-group" data-toggle="buttons">
									<?php foreach($CartConfig['product']['cycle'] as $cycle_key=>$cycle): ?>
									<label class="btn btn-custom ">
										<input type="radio" name="billingcycle" value="<?php echo $cycle['billingcycle']; ?>" <?php if($CartConfig['billingcyle']==$cycle['billingcycle']): ?>checked="" <?php elseif($cycle['is_checked']==1 &&
											!$CartConfig['billingcyle'] && !$addParam['billingcycle']): ?> checked="" <?php elseif($addParam['billingcycle']==$cycle['billingcycle'] &&
										!$CartConfig['billingcyle']): ?> checked="" <?php endif; ?>><?php echo $cycle['billingcycle_zh']; if($cycle['cycle_discount']): ?>
										<span class="cycle-select-discount"><?php echo $cycle['cycle_discount']; ?><?php echo $Lang['fracture']; ?></span>
										<?php endif; ?>
									</label>
									<?php endforeach; ?>
								</div>
							</div>
						</div>

						<div class="form-group row" <?php if($CartConfig['product']['host']['show']==0): ?> style="display:none" <?php endif; ?>>
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $Lang['host_name']; ?></label>
							<div class="col-md-3">
								<input class="form-control" type="text" disabled="disabled" <?php if($CartConfig['host']): ?>value="<?php echo $CartConfig['host']; ?>" <?php else: ?>value="<?php echo $CartConfig['product']['host']['host']; ?>"
									<?php endif; ?> />
								<input class="form-control" name="host" type="hidden" <?php if($CartConfig['host']): ?>value="<?php echo $CartConfig['host']; ?>" <?php else: ?>value="<?php echo $CartConfig['product']['host']['host']; ?>"
									<?php endif; ?> />
							</div>
						</div>
						<div class="form-group row" <?php if($CartConfig['product']['password']['show']==0): ?> style="display:none" <?php endif; ?>>
							<label for="example-search-input" class="col-md-2 col-form-label"><?php echo $Lang['host_password']; ?></label>
							<div class="col-8 col-md-3">
								<input
								class="form-control getPassword"
								id="password"
								name="password"
								<?php if($CartConfig['password']): ?>
									value="<?php echo $CartConfig['password']; ?>"
								<?php else: ?>
									 value="<?php echo $CartConfig['product']['password']['password']; ?>"
								<?php endif; ?>
								/>
							</div>
							<div class="d-inline-flex align-items-center">
								<span class="bx bx-show-alt pointer fs-20 text-primary" onclick="showPwd()"></span>
							</div>
							<div class="col-1 col-md-1 fs-18 d-flex align-items-center">

								<i class="fas fa-dice create_random_pass pointer" onclick="randomPwd()"></i>
							</div>
							<label id="password-error-tip" class="control-label error-tip" for="password"></label>
						</div>



					</div>
				</div>
			</div>
			<div class="col-xl-4">

				<div class="card">
					<div class="card-body">
						<h4 class="card-title fs-16"><?php echo $Lang['order_summary']; ?>：</h4>

						<div class="table-responsive configoption_total fs-14">

						</div>
					</div>
				</div>
				<!-- end card -->
			</div>
		</div>
	</form>
	<!-- select -->
	<link rel="stylesheet"
		href="/themes/cart/default/assets/js/bootstrap-select/css/bootstrap-select.min.css?v=<?php echo $Ver; ?>">
	<script src="/themes/cart/default/assets/js/bootstrap-select/js/bootstrap-select.min.js?v=<?php echo $Ver; ?>"></script>
	<!-- slider -->
	<link rel="stylesheet"
		href="/themes/cart/default/assets/js/ion-rangeslider/css/ion.rangeSlider.min.css?v=<?php echo $Ver; ?>">
	<script src="/themes/cart/default/assets/js/ion-rangeslider/js/ion.rangeSlider.min.js?v=<?php echo $Ver; ?>"></script>
	<script>
		// 密码重置
		var pwdRule = {
			len_num: '<?php echo $CartConfig['product']['password']['rule']['len_num']; ?>',
			num: '<?php echo $CartConfig['product']['password']['rule']['num']; ?>',
			upper: '<?php echo $CartConfig['product']['password']['rule']['upper']; ?>',
			lower: '<?php echo $CartConfig['product']['password']['rule']['lower']; ?>',
			special: '<?php echo $CartConfig['product']['password']['rule']['special']; ?>'
		}
		var links = <?php echo json_encode($CartConfig['links']); ?>;
		var order_frm_tpl = '<?php echo $order_frm_tpl; ?>';
		var tpl_type = '<?php echo $tpl_type; ?>';

		$(document).ready(function(){
			let dom=$('.dropdown-menu .inner .show')
			for(let item of dom){
				item.style.maxHeight="300px"
			}
		});

		$(function(){
			$("#addCartForm").on('blur',".getVal",function(e){
				let arr=e.target.dataset.object.split('||')
				let reg=''
				if(arr[1]!='') reg=new RegExp(arr[1])
				if(arr[0]==1){
					if(e.target.value==''){
						e.target.classList.add("is-invalid");
					}else if(reg!='' && reg.test(e.target.value)==false){
						e.target.classList.add("is-invalid");
					}else{
						e.target.classList.remove("is-invalid");
					}
				}
			})
		})

		$(function(){
			$("#addCartForm").on('blur',".getTextareaVal",function(e){

				let arr=e.target.dataset.object.split('||')
				let reg=''
				if(arr[1]!='') reg=new RegExp(arr[1])
				if(arr[0]==1){
					if(e.target.value==''){
						e.target.classList.add("is-invalid");
					}else if(reg!='' && reg.test(e.target.value)==false){
						e.target.classList.add("is-invalid");
					}else{
						e.target.classList.remove("is-invalid");
					}
				}
			})
		})

		function randomPwd () {
			verConfigGPsd()
		}

	// 获取主机密码校验规则
	var passwordRules =  <?php echo json_encode($CartConfig['product']['password']['rule']); ?>;
	var showPassword = <?php echo json_encode($CartConfig['product']['password']['show']); ?>;

	$(document).on('blur', '.getPassword', function(){
		verConfigGPsd()
	})

	function verConfigGPsd() {
		setTimeout(function(){
			let result = checkingPwd1($(".getPassword").val(), passwordRules.num, passwordRules.upper, passwordRules.lower, passwordRules.special)
			if(result.flag) {
				$('#password-error-tip').css('display','none');
				$('.getPassword').removeClass("is-invalid");
			}else{
				$("#password-error-tip").html(result.msg);
				$(".getPassword").addClass("is-invalid");
				$('#password-error-tip').css('display','block');
			}
		},10)
	}
	</script>
	<script src="/themes/cart/default/assets/js/configureproduct.js?v=<?php echo $Ver; ?>"></script>
	<!-- 滑块禁用区域 -->
	<script>
		$('.configoption_range').each(function(){
			let sub = $(this).data('sub');
			let max = parseFloat($(this).attr('max'));
			let min = parseFloat($(this).attr('min'));
			let inputWidth = parseFloat($(this).width());
			let oneWidth = inputWidth / (max - min)
			let keyArr = [];
			sub.map(item => {
				let itemMin = parseFloat(item.qty_minimum)
				let itemMax = parseFloat(item.qty_maximum)
				keyArr.push(itemMin)
				for(var i = itemMin; i<itemMax; i++) {
					keyArr.push(i)
				}
				keyArr.push(itemMax)
			})
			keyArr = Array.from(new Set(keyArr));
			for(var t = min; t<max; t++) {
				if(keyArr.indexOf(t) == -1) {
					$(this).after('<span class="range_none" title="'+t+'" style="width: '+oneWidth+'px; left: '+oneWidth*(t-min)+'px "></span>')
				}
			}
		})
	</script>
<?php if($TplName != 'login' && $TplName != 'register' && $TplName != 'pwreset' && $TplName != 'bind' && $TplName != 'loginaccesstoken'): ?>
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