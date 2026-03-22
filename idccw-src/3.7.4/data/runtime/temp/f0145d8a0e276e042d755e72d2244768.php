<?php /*a:1:{s:105:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/nat.tpl";i:1737452904;}*/ ?>
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