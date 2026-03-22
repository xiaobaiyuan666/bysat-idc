<?php /*a:1:{s:123:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/servicedetail-billing.tpl";i:1737452904;}*/ ?>
<div class="table-responsive">
  <table class="table mb-0 mt-3">
    <thead class="thead-light">
      <tr>
        <th><?php echo $Lang['payment_time']; ?></th>
        <th><?php echo $Lang['source']; ?></th>
        <th><?php echo $Lang['payment_amount']; ?></th>
        <th><?php echo $Lang['serial_number']; ?></th>
        <th><?php echo $Lang['payment_method']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($HostRecharge as $item): ?>
      <tr>
        <td><?php echo date("Y-m-d H:i:s",!is_numeric($item['pay_time'])? strtotime($item['pay_time']) : $item['pay_time']); ?></td>
        <td><?php echo $item['type']; ?></td>
        <td><?php echo $Currency['prefix']; ?><?php echo $item['amount_in']; ?><?php echo $Currency['suffix']; ?></td>
        <td><?php echo $item['trans_id']; ?></td>
        <td><?php echo $item['gateway']; ?></td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
<div class="table-footer">
  <div class="table-tools">

  </div>
  <div class="table-pagination">
    <div class="table-pageinfo mr-2">
      <span><?php echo $Lang['common']; ?> <?php echo $Total; ?> <?php echo $Lang['strips']; ?></span>
      <span class="mx-2">
        <?php echo $Lang['each_page']; ?>
        <select name="" id="limitSel" class="billing-limit" onchange="getGeneralBilling('<?php echo app('request')->get('id'); ?>', '<?php echo app('request')->get('page'); ?>')">
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