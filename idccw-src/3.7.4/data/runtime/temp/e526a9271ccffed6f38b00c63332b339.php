<?php /*a:1:{s:112:"C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4/public/themes/clientarea/default/includes/clientarea-list.tpl";i:1737452904;}*/ ?>
<!-- start：资源列表 -->
<h4 class="card-title mt-0"><?php echo $Lang['resource_list']; ?></h4>
<div class="h100p user-center_resources table-responsive">
  <table class="table tablelist mb-0 mt-3">
		<colgroup>
			<col width="25%">
			<col width="25%">
			<col width="25%">
			<col width="15%">
			<col width="10%">
		</colgroup>
    <thead class="bg-light">
      <tr>
        <th><?php echo $Lang['machine_status']; ?></th>
        <th><?php echo $Lang['host_name']; ?></th>
        <th class="pointer" prop="nextduedate"><?php echo $Lang['due_date']; ?>
          <span class="text-black-50 d-inline-flex flex-column justify-content-center ml-1 offset-3">
            <i class="bx bx-caret-up"></i>
            <i class="bx bx-caret-down"></i>
          </span>
        </th>
        <th><?php echo $Lang['cost']; ?></th>
        <th>IP</th>
      </tr>
    </thead>
    <tbody>
      <?php if($ClientArea['hostlist']): foreach($ClientArea['hostlist'] as $list): ?>
      <tr>
        <td scope="row">
          <span class="user-center_dot mr-1 bg-success"></span>
          <?php echo $list['domainstatus_desc']; ?>
        </td>
        <td>
          <a href="servicedetail?id=<?php echo $list['id']; ?>" class="text-dark"><?php echo $list['productname']; ?>(<?php echo $list['domain']; ?>)</a>
        </td>
        <td><?php if($list['billingcycle']!="free" && $list['cycle_desc']!='一次性'): ?><?php echo date("Y-m-d H:i",!is_numeric($list['nextduedate'])? strtotime($list['nextduedate']) : $list['nextduedate']); else: ?> - <?php endif; ?></td>
        <td><?php if($list['billingcycle']!="free"): ?><?php echo $list['price_desc']; ?>/<?php echo $list['cycle_desc']; else: ?><?php echo $list['cycle_desc']; ?><?php endif; ?></td>
        <td><?php echo $list['dedicatedip']; ?></td>
      </tr>
      <?php endforeach; else: ?>
      <tr>
        <td colspan="8">
          <div class="no-data"><?php echo $Lang['nothing_content']; ?></div>
        </td>
      </tr>
      <?php endif; ?>
    </tbody>
  </table>
</div>
<!-- 表单底部调用开始 -->
<div class="table-footer">
  <div class="table-tools">

  </div>
  <div class="table-pagination">
    <div class="table-pageinfo mr-2">
      <span><?php echo $Lang['common']; ?> <?php echo $ClientArea['Total']; ?> <?php echo $Lang['strips']; ?></span>
      <span class="mx-2">
        <?php echo $Lang['each_page']; ?>
        <select name="" id="sourcelimitSel">
          <option value="5" <?php if($ClientArea['Limit']==5): ?>selected<?php endif; ?>>5 </option>
          <option value="10" <?php if($ClientArea['Limit']==10): ?>selected<?php endif; ?>>10 </option>
          <option value="15" <?php if($ClientArea['Limit']==15): ?>selected<?php endif; ?>>15 </option>
          <option value="20" <?php if($ClientArea['Limit']==20): ?>selected<?php endif; ?>>20 </option>
          <option value="50" <?php if($ClientArea['Limit']==50): ?>selected<?php endif; ?>>50 </option>
          <option value="100" <?php if($ClientArea['Limit']==100): ?>selected<?php endif; ?>>100
          </option>
        </select> <?php echo $Lang['strips']; ?></span>
    </div>
    <ul class="pagination pagination-sm">
      <?php echo $ClientArea['Pages']; ?>
    </ul>
  </div>
</div>
<script>
  $(function () {
    // 排序
    $('.bg-light th').on('click', function () {
      if (!$(this).attr('prop')) return false
      getSourceList('sort', $(this).attr('prop'))
       
    })
    //改变排序样式
    changeStyle()
    function changeStyle() {
        $('.bg-light th.pointer span').children().css('color','rgba(0, 0, 0, 0.1)')
        let index
        let sort = localStorage.getItem('sort')
        if( sort=== 'desc') {
          index = 1
        } else if(sort=== 'asc'){
          index = 0
        }
        $('.bg-light th.pointer span').children().eq(index).css('color','rgba(0, 0, 0, 0.8)')
    }
  

    // 每页数量选择改变
    $('#sourcelimitSel').on('change', function () {
      getSourceList('limit')
    })

    // 分页
    $('.page-link').on('click', function (e) {
      e.preventDefault()
      $.get('' + $(this).attr('href'), function (data) {
        $('#sourceListBox').html(data)
      })
    })
  })

  function getSourceList(searchType, orderby) {
    // 搜索条件
    var searchObj = {
      action: 'list'
    }

    if (searchType == 'sort') {
      searchObj.sort = (localStorage.getItem('sort') == 'asc') ? 'desc' : 'asc'
      localStorage.setItem('sort', searchObj.sort)
      searchObj.orderby = orderby
       }
    if (searchType == 'limit') {
      searchObj.limit = $('#sourcelimitSel').val()
      searchObj.page = 1
    }
    $.ajax({
      type: "get",
      url: '' + '/clientarea',
      data: searchObj,
      success: function (data) {
        $('#sourceListBox').html(data)
      }
    });
  }
</script>
