<?php /*a:5:{s:98:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/clientarea.tpl";i:1737452904;s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/head.tpl";i:1737452904;s:101:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/menu.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/pageheader.tpl";i:1737452904;s:107:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/breadcrumb.tpl";i:1737452904;}*/ ?>

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
				<?php endif; ?>
<script src="/themes/clientarea/default/assets/libs/echarts/echarts.min.js?v=<?php echo $Ver; ?>"></script>
<script>
  $(function () {
    var credit = '<?php echo $ClientArea['index']['client']['credit']; ?>' // 余额
    var unpaid = '<?php echo $ClientArea['index']['invoice_unpaid']; ?>' // 未支付
    var creditNum = parseFloat(credit.replace('<?php echo $Lang['element']; ?>', ''))
    var unpaidNum = parseFloat(unpaid.replace('<?php echo $Lang['element']; ?>', ''))

    var percentage = parseFloat(creditNum / (creditNum + unpaidNum)).toFixed(2)

    if (creditNum === 0 || (creditNum === 0 && unpaidNum === 0)) {
      percentage = 0
    }
    var myChart = echarts.init(document.getElementById("balanceCharts"));
    var
      option = {
        title: {
          text: '{a|' + credit + '}\n\n{c|' + '<?php echo $Lang['current_balance']; ?>' + '}',
          x: 'center',
          y: 'center',
          textStyle: {
            rich: {
              a: {
                fontSize: 16,
                color: '#007bfc'
              },
              c: {
                fontSize: 12,
                color: '#000000',
                padding: [5, 0]
              },

            }
          }
        },
        series: [
          {
            type: 'gauge',
            radius: '80%',
            clockwise: true,
            startAngle: '270',
            endAngle: '-89.9999',
            //调整间隔距离
            splitNumber: 0,
            detail: {
              offsetCenter: [0, -20],
              formatter: ' '
            },
            pointer: {
              show: false
            },
            axisLine: {
              show: true,
              lineStyle: {
                color: [
                  [0, '#228cfc'],
                  //计算比例
                  [percentage, '#228cfc'],

                  [1, '#efefef']
                ],
                width: 9
              }
            },
            axisTick: {
              show: false
            },
            splitLine: {
              show: false,
              length: 32,
              lineStyle: {
                color: '#fff',
                width: 6
              }
            },
            axisLabel: {
              show: false
            }
          }
        ]
      };
    myChart.setOption(option, true);
    
    // echarts 随屏幕大小改变大小
    window.addEventListener("resize",function () {
      myChart.resize();
    });

    // 资源列表
    getSourceList()
  })

  function getSourceList() {
    $('#sourceListBox').html(`<div class="h-100 d-flex align-items-center justify-content-center"><?php echo $Lang['data_loading']; ?>......</div>`)
    $.ajax({
      type: "get",
      url: '/clientarea',
      data: {
        action: 'list'
      },
      success: function (data) {
        $('#sourceListBox').html(data)
      }
    });
  }

</script>

<div class="row">
  <!-- start：个人信息 -->
  <section class="col-md-12 col-xl-4">
    <div class="card card-body user-center_h300 p-0 pb-3">
      <!-- old -->
      <!-- <div class="d-flex align-items-center justify-content-center mb-4">
        <div class="mr-4 d-flex align-items-center justify-content-center user-center_header"><?php if(preg_match("/[\x7f-\xff]/", substr($Userinfo['user']['username'],0,3))): ?> <?php echo substr($Userinfo['user']['username'],0,3); else: ?>
          <?php echo strtoupper(substr($Userinfo['user']['username'],0,1)); ?> <?php endif; ?></div>
        <div class="ml-20">
          <p>
            <span class="user-center_name"><?php echo $Userinfo['user']['username']; ?></span>
            <span>ID:<?php echo $Userinfo['user']['id']; ?></span>
          </p>
          <span><?php if($Userinfo['user']['phonenumber']): ?><?php echo substr($Userinfo['user']['phonenumber'],0,3); ?>***<?php echo substr($Userinfo['user']['phonenumber'],7); ?><?php endif; ?></span>
        </div>
      </div>
      <div class="d-flex justify-content-around user-center_safety_wrapper mb-3">
        <a href="security" class="user-center_safety" <?php if($Userinfo['user']['certifi']['status']!=1): ?>
          style="background-color:#aaa;" <?php endif; ?>>
          <i class="bx bx-check-shield"></i>
        </a>
        <a href="security" class="user-center_safety" <?php if(!$Userinfo['user']['phonenumber']): ?> style="background-color:#aaa;"
          <?php endif; ?>>
          <i class="bx bx-phone-call"></i>
        </a>
        <a href="security" class="user-center_safety" <?php if(!$Userinfo['user']['email']): ?> style="background-color:#aaa;" <?php endif; ?>>
          <i class="bx bx-mail-send"></i>
        </a>
      </div> -->
      <!-- old -->
      <div class="h-75">
        <div class="h-75 d-flex align-items-center justify-content-around flex-column" style="margin-top:20px">
          <div class=" d-flex align-items-center justify-content-center" style="">
            <div
              style="background-color: #fff;position: relative;left:-15px;top: 8px; width: 62px;height: 62px;border-radius: 50%;padding: 3px">
              <div class="mr-4 d-flex align-items-center justify-content-center user-center_header">
                <?php if(preg_match("/^[0-9]*[A-Za-z]+$/is", substr($Userinfo['user']['username'],0,1))): ?> 
                  <?php echo strtoupper(substr($Userinfo['user']['username'],0,1)); elseif(preg_match("/^[\x7f-\xff]*$/", substr($Userinfo['user']['username'],0,3))): ?> 
                  <?php echo substr($Userinfo['user']['username'],0,3); else: ?>
                  <?php echo strtoupper(substr($Userinfo['user']['username'],0,1)); ?> 
                <?php endif; ?>
              </div>
            </div>
            <div class="ml-20">
              <div>
                <span class="user-center_name"><?php echo $Userinfo['user']['username']; ?></span>
                <span>ID:<?php echo $Userinfo['user']['id']; ?></span>
              </div>
              <span><?php if($Userinfo['user']['phonenumber']): ?><?php echo substr($Userinfo['user']['phonenumber'],0,3); ?>***<?php echo substr($Userinfo['user']['phonenumber'],7); ?><?php endif; ?></span>
            </div>
          </div>
          <div class="d-flex align-items-center justify-content-around user-center_safety_wrapper" style="margin-bottom:-25px">
            <a href="security" class="user-center_safety" <?php if($Userinfo['user']['certifi']['status']!=1): ?>
              style="background-color:#aaa;" <?php endif; ?>>
              <i class="bx bx-check-shield"></i>
            </a>
            <a href="security" class="user-center_safety" <?php if(!$Userinfo['user']['phonenumber']): ?>
              style="background-color:#aaa;" <?php endif; ?>>
              <i class="bx bx-phone-call"></i>
            </a>
            <a href="security" class="user-center_safety" <?php if(!$Userinfo['user']['email']): ?> style="background-color:#aaa;"
              <?php endif; ?>>
              <i class="bx bx-mail-send"></i>
            </a>
          </div> 
        </div>
      </div>
      <hr>
      <div class="d-flex justify-content-around align-items-center h-25">
        <a href="supporttickets" class="text-dark">
          <div class="text-center text-warning"><?php echo $ClientArea['index']['ticket_count']; ?></div>
          <span><?php echo $Lang['pending_work_order']; ?></span>
        </a>
        <a href="billing" class="text-dark">
          <div class="text-center text-primary"><?php echo $ClientArea['index']['order_count']; ?></div>
          <span><?php echo $Lang['unpaid_order']; ?></span>
        </a>
        <div>
          <div class="text-center text-danger"><?php echo $ClientArea['index']['host']; ?></div>
          <span><?php echo $Lang['number_of_products']; ?></span>
        </div>
      </div>
    </div>
  </section>

  <!-- start：财务信息 -->
  <section class="col-md-12 col-xl-4">
    <div class="card card-body user-center_h300">
      <div class="d-flex h100p">
        <div class="d-flex flex-column align-items-center justify-content-center flex1">
          <div class="w-100 h-100 d-flex justify-content-center" id="balanceCharts"></div>
        </div>
        <div class="d-flex flex-column justify-content-center flex1">
          <span class="d-inline-flex fz-12 text-black-50 mb-2"><?php echo $Lang['records_of_consumption']; ?></span>
          <span class="d-inline-flex fz-14 text-black-80"><?php echo $Lang['consumption_this_month']; ?>：<?php echo $ClientArea['index']['intotal']; ?></span>
          <span class="d-inline-flex fz-12 text-black-50 mt-4 mb-2"><?php echo $Lang['order_records']; ?></span>
          <span class="d-inline-flex fz-14 text-black-80"><?php echo $Lang['unpaid']; ?>：<?php echo $ClientArea['index']['invoice_unpaid']; ?></span>
          <?php if($ClientArea['index']['allow_recharge'] == '1'): ?>
          <a href="/addfunds" class="btn btn-primary mt-5 w-50"><?php echo $Lang['recharge']; ?></a>
          <?php endif; ?>
        </div>
      </div>
    </div>
  </section>
  <!-- start：已开通产品 -->
  <section class="col-md-12 col-xl-4">
    <div class="card card-body user-center_h300">
      <h4 class="card-title mt-0"><?php echo $Lang['products_launched_all']; ?></h4>
      <div class="user-center_product_grid mt-3" style="overflow-y: auto;">
        <?php foreach($ClientArea['index']['host_nav'] as $list): ?>
        <a href="service?groupid=<?php echo $list['id']; ?>" class="user-center_product">
          <span>
            <i class="bx bxs-grid-alt"></i>
            <?php echo $list['groupname']; ?>
          </span>
          <span>(<?php echo $list['count']; ?>)</span>
        </a>
        <?php endforeach; ?>
      </div>
    </div>
  </section>


  <!-- start：资源列表 -->
  <section class="col-md-12 col-xl-8">
    <div class="card card-body user-center_calc" id="sourceListBox">

    </div>
  </section>


  <!-- start：公告通知 -->
  <section class="col-md-12 col-xl-4">
    <div class="card card-body user-center_calc">
      <h4 class="mb-4 card-title d-flex justify-content-between">
        <span><?php echo $Lang['announcement']; ?></span>
        <a href="news" class="fs-12 font-weight-normal"><?php echo $Lang['view_more']; ?></a>
      </h4>
      <div class="user-center_notice h100p">
        <ul class="user-center_notice_ul pl-0">
          <?php if($ClientArea['index']['news']): foreach($ClientArea['index']['news'] as $list): ?>
          <li class="user-center_notice_item">
            <span class="notice_item_time text-black-50"><?php echo date("Y-m-d H:i",!is_numeric($list['push_time'])? strtotime($list['push_time']) : $list['push_time']); ?></span>
            <a href="newsview?id=<?php echo $list['id']; ?>" class="notice_item_title"><?php echo $list['title']; ?></a>
          </li>
          <?php endforeach; else: ?>
          <tr>
            <td colspan="2">
              <div class="no-data"><?php echo $Lang['nothing']; ?></div>
            </td>
          </tr>
          <?php endif; ?>
        </ul>
      </div>
    </div>
  </section>
</div><?php if($TplName != 'login' && $TplName != 'register' && $TplName != 'pwreset' && $TplName != 'bind' && $TplName != 'loginaccesstoken'): ?>
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