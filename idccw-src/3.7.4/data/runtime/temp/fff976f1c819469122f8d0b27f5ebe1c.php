<?php /*a:1:{s:119:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/servicedetail/servicedetail-log.tpl";i:1737452904;}*/ ?>

<div class="table-responsive">
  <table class="table mb-0 mt-3">
    <thead class="thead-light">
      <tr>
        <th><?php echo $Lang['operation_time']; ?></th>
        <th><?php echo $Lang['operation_details']; ?></th>
        <th><?php echo $Lang['operator']; ?></th>
        <th><?php echo $Lang['ip_address']; ?></th>
      </tr>
    </thead>
    <tbody>
      <?php foreach($RecordLog as $item): ?>
      <tr>
        <td><?php echo date("Y-m-d H:i:s",!is_numeric($item['create_time'])? strtotime($item['create_time']) : $item['create_time']); ?></td>
        <td><?php echo $item['description']; ?></td>
        <td><?php echo $item['user']; ?></td>
        <td><?php echo $item['ipaddr']; ?></td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
</div>
<div class="table-footer">
  <div class="table-tools">

  </div>
  <div class="table-pagination generallog">
    <div class="table-pageinfo mr-2">
      <span><?php echo $Lang['common']; ?> <?php echo $Total; ?> <?php echo $Lang['strips']; ?></span>
      <span class="mx-2">
        <?php echo $Lang['each_page']; ?>
        <select name="" id="limitSel" class="log-limit" onchange="getGeneralLog('<?php echo app('request')->get('id'); ?>', '<?php echo app('request')->get('page'); ?>')">
          <option value="10" <?php if($Limit==10): ?>selected<?php endif; ?>>10 </option>
          <option value="15" <?php if($Limit==15): ?>selected<?php endif; ?>>15 </option>
          <option value="20" <?php if($Limit==20): ?>selected<?php endif; ?>>20 </option>
          <option value="50" <?php if($Limit==50): ?>selected<?php endif; ?>>50 </option>
          <option value="100" <?php if($Limit==100): ?>selected<?php endif; ?>>100
          </option>
        </select> <?php echo $Lang['strips']; ?> </span>
    </div>
    <ul class="pagination pagination-sm">
      <?php echo $Pages; ?>
    </ul>
  </div>
</div>

<script>

  $(document).on('click', ".generallog a[class='page-link']", function () {
    var _this = $(this)
    $.ajax({
      url: _this.prop('href'),
      type: 'GET',
      success: function (res) {
        $('#settings1').html(res)
      }
    })

    return false;
  });
</script>