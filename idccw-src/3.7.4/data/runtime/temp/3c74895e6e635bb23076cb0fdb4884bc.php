<?php /*a:1:{s:94:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/cart/default/ordersummary.tpl";i:1737452904;}*/ ?>


<style>
	.ordersummary td {
		border:none!important;
		padding: 5px!important;
	}
	.mobile-bottom-total{
		display: none;
		position: fixed;
		width: 100%;
		height: 4rem;
		text-align: right;
		line-height: 4rem;
		padding-right: 1.25rem;
		background: #FFFFFF;
		box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
		z-index: 99999;
		bottom: 0;
		left: 0;
		-webkit-overflow-scrolling: auto;
	}
	.mobile-bottom-total .total {
		font-size: 18px;
		font-weight: 400;
		color: #333333;
		margin-right: 0.5rem;
	}
	.btn.active { 
		z-index: 0 !important; 
	}
</style>

<table class="table mb-5 mt-2 ordersummary">
	<tbody>
		<tr>
			<td class="color-999 mw-250"><?php echo $ConfigureTotal['product_name']; ?>:</td>
			<td class="text-right"><?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo $ConfigureTotal['product_price']; if($ConfigureTotal['product_setup_fee']>0): ?>
					+<?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo $ConfigureTotal['product_setup_fee']; ?><?php echo $Lang['initial_installation_fee']; ?>
				<?php endif; ?>
			</td>
		</tr>
		<?php foreach($ConfigureTotal['child'] as $configure): ?>
		<tr>
			<td class="mw-250">
				<span class="color-999"><?php echo $configure['option_name']; ?>:</span>
				<?php if($configure['option_type'] == '12'): if($configure['icon_flag']): ?>
						<img class='mr-1' src='/upload/common/country/<?php echo $configure['icon_flag']; ?>.png' height='15'/>
					<?php endif; elseif($configure['option_type'] == '5'): if($configure['icon_os']): ?>
						<img class='mr-1' src='/upload/common/system/<?php echo $configure['icon_os']; ?>.svg' height='20'/>
					<?php endif; ?>
				<?php endif; ?>
				<span style="word-break: break-all;white-space: normal;">
					<?php if($configure['qty']): ?><?php echo $configure['qty']; else: ?><?php echo $configure['sub_name']; ?><?php endif; ?>
				</span>
			</td>
			<td class="text-right"><?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo $configure['suboption_price']; if($configure['suboption_setup_fee'] > 0): ?> + <?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo $configure['suboption_setup_fee']; ?><?php echo $Lang['initial_installation_fee']; ?><?php endif; ?></td>
		</tr>
		<?php endforeach; ?>

		<tr>
			<td class="pr-0">
				<hr class="my-2">
			</td>
			<td class="pl-0">
				<hr class="my-2">
			</td>
		</tr>
		<!-- 有折扣时， 才显示总价 -->
		<?php if($ConfigureTotal['type']): ?>
			<tr>
				<td class="color-999"><?php echo $Lang['price']; ?>:</td>
				<td class="text-dark text-right">
					<?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo $ConfigureTotal['total']; ?>
				</td>
			</tr>
		<?php endif; if($ConfigureTotal['type']): ?>
			<tr>
				<td class="color-999">
					<?php if($ConfigureTotal['type']['type'] == '1'): ?>
						<span class=""><?php echo $Lang['customer_discount_price']; ?>
							(<span class="discount-num"></span><?php echo $Lang['fracture']; ?>):
						</span>
					<?php elseif($ConfigureTotal['type']['type'] == '2'): ?>
						<span class=""><?php echo $Lang['customer_discount_province']; ?>
							<span class="discount-num"><?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo $ConfigureTotal['type']['bates']; ?>):</span>
						</span>
					<?php endif; ?>
				</td>
				<td class="text-dark text-right">-
					<?php echo $ConfigureTotal['currency']['prefix']; ?>
					<?php echo bcsub(bcsub($ConfigureTotal['total'],$ConfigureTotal['sale_setupfee_total']),$ConfigureTotal['sale_signal_price']); ?>
				</td>
			</tr>

		<?php endif; ?>

		<tr class="mobile-hide">
			<td class="color-999"><?php echo $Lang['total_price']; ?>:</td>
			<td class="font-weight-bold text-dark text-right">
				<?php if(!$ConfigureTotal['type']): ?>
					<?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo bcadd($ConfigureTotal['signal_price'],$ConfigureTotal['signal_setupfee']); else: ?>
					<?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo bcadd($ConfigureTotal['sale_signal_price'],$ConfigureTotal['sale_setupfee_total']); ?>
				<?php endif; ?>
			</td>
		</tr>
	</tbody>
</table>

<div class="text-sm-right mt-4 mt-sm-0 mobile-hide" style="position: absolute;bottom: 20px;right: 20px;">
							<button type="button" style="cursor: pointer;" class="btn btn-primary" id="addToCartBtn"><i class="mdi mdi-cart-arrow-right mr-1"></i><?php echo $Lang['add_cart']; ?></button>
						</div>

	<!-- 移动端底部价格展示 -->
	<div class="mobile-bottom-total">
     	<span class="total"><?php echo $Lang['total_price']; ?>：	<?php if(!$ConfigureTotal['type']): ?>
					<?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo bcadd($ConfigureTotal['signal_price'],$ConfigureTotal['signal_setupfee']); else: ?>
					<?php echo $ConfigureTotal['currency']['prefix']; ?><?php echo bcadd($ConfigureTotal['sale_signal_price'],$ConfigureTotal['sale_setupfee_total']); ?>
				<?php endif; ?></span>
			 <button type="button" style="cursor: pointer;" class="btn btn-primary" id="addToCartBtnTwo"><i class="mdi mdi-cart-arrow-right mr-1"></i><?php echo $Lang['add_cart']; ?></button>
		</div>
<script>
	$(function() {
		if(navigator.userAgent.match(/mobile/i)) {
			$('.mobile-bottom-total').show()
			$('.mobile-hide').remove()
		}else{
			$('.mobile-bottom-total').remove()
		}
		console.log('2',$(".getPassword"))
		console.log('val', $(".getPassword").val())
		console.log('showPassword2',showPassword)
		// 产品信息
		var products = <?php echo json_encode($ConfigureTotal); ?>;
		// 订单折扣量
		if ('<?php echo $ConfigureTotal['type']['type']; ?>' == '1') {
			if (parseFloat(products.type.bates) % 10 == 0) {
				$('.discount-num').text(parseFloat(parseFloat(products.type.bates) / 10))
			} else {
				$('.discount-num').text(parseFloat(parseFloat(products.type.bates) / 10).toFixed(2))
			}
		}

    if (!!window.ActiveXObject || "ActiveXObject" in window){
      console.log('ie');
      $('#addToCartBtn,#addToCartBtnTwo').click(function () {
		let result = {flag:true}
		if(passwordRules != null && showPassword == 1) {
			result = checkingPwd1($(".getPassword").val(), passwordRules.num, passwordRules.upper, passwordRules.lower, passwordRules.special)
		}
        if(result.flag) {
						$('#addCartForm').submit()
					}
					else {
						toastr.error($('.is-invalid').parents('.form-group').find('.error-tip').html())
					}
				// $('#addCartForm').submit()

		  })
    }else{
			// console.log('不是ie');
      $('#addToCartBtn,#addToCartBtnTwo').click(function () {
		let result = {flag:true}
		if(passwordRules != null && showPassword == 1) {
			result = checkingPwd1($(".getPassword").val(), passwordRules.num, passwordRules.upper, passwordRules.lower, passwordRules.special)
		}
					if(result.flag) {
						$('#addCartForm').submit()
					}
					else {
						toastr.error($('.is-invalid').parents('.form-group').find('.error-tip').html())
					}
					// $('#addCartForm').submit()
       })
      // 加入购物车 提交按钮
		// $(document).on('click', '#addToCartBtn', function () {
		// 	if($(this).data("disabled")) return false;
		// 	var InpCheck = $( "input[name^='customfield']")
		// 	var textareaCheck = $( "textarea[name^='customfield']")
		// 	if (checkListFormVerify([...InpCheck, ...textareaCheck])){
		// 		var position = $(".is-invalid:first").offset();
		// 		scrolltop = position.top-70;
		// 		$("html,body").animate({scrollTop:scrolltop}, 1000);
		// 		return false;
		// 	}
		// 	$(this).prepend('<i class="bx bx-loader bx-spin font-size-16 align-middle mr-2"></i>')
		// 	$(this).data("disabled",true);
    //   // debugger
		// 	$('#addCartForm').submit()
		// })
    }
   
		
	})
</script>

