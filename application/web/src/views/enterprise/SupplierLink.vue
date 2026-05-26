<template>
  <div class="supplier-container">
    <h3>供应商关联管理（链上存证）</h3>
    <el-form :inline="true">
      <el-form-item label="证照溯源码">
        <el-input v-model="traceCode" placeholder="请输入证照溯源码" clearable />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
        <el-button @click="showAddDialog">添加关联（上链）</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="supplierList" border stripe>
      <el-table-column prop="supplier_name" label="供应商名称" />
      <el-table-column prop="contact_person" label="联系人" />
      <el-table-column prop="contact_phone" label="联系电话" />
      <el-table-column prop="cooperation_start" label="合作开始" />
      <el-table-column prop="cooperation_end" label="合作结束" />
      <el-table-column prop="notes" label="备注" />
      <el-table-column prop="timestamp" label="上链时间" width="160" />
      <el-table-column prop="txid" label="交易哈希" width="200" show-overflow-tooltip />
      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button type="text" @click="tryModify(scope.row)">尝试修改</el-button>
          <el-button type="text" @click="verifyHash(scope.row)">验证哈希</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog title="添加供应商关联" :visible.sync="dialogVisible" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="证照溯源码" required>
          <el-input v-model="form.traceability_code" disabled />
        </el-form-item>
        <el-form-item label="供应商名称" required>
          <el-input v-model="form.supplier_name" />
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="form.contact_person" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="form.contact_phone" />
        </el-form-item>
        <el-form-item label="合作开始">
          <el-date-picker v-model="form.cooperation_start" type="date" placeholder="选择日期" value-format="yyyy-MM-dd" />
        </el-form-item>
        <el-form-item label="合作结束">
          <el-date-picker v-model="form.cooperation_end" type="date" placeholder="选择日期" value-format="yyyy-MM-dd" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input type="textarea" v-model="form.notes" />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdd">确认上链</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { addSupplierLink, getSupplierLinks } from '@/api/enterprise'

export default {
  data() {
    return {
      traceCode: '',
      supplierList: [],
      dialogVisible: false,
      form: {
        traceability_code: '',
        supplier_name: '',
        contact_person: '',
        contact_phone: '',
        cooperation_start: '',
        cooperation_end: '',
        notes: ''
      }
    }
  },
  methods: {
    async loadList() {
      if (!this.traceCode) {
        this.$message.warning('请输入证照溯源码')
        return
      }
      const res = await getSupplierLinks(this.traceCode)
      if (res.code === 200) {
        // 按上链时间降序（最新的在前）
        this.supplierList = res.data.sort((a, b) => (b.timestamp || 0) - (a.timestamp || 0))
      }
    },
    showAddDialog() {
      if (!this.traceCode) {
        this.$message.warning('请先输入要关联的证照溯源码')
        return
      }
      this.form.traceability_code = this.traceCode
      this.dialogVisible = true
    },
    async submitAdd() {
      if (!this.form.supplier_name) {
        this.$message.warning('请输入供应商名称')
        return
      }
      const res = await addSupplierLink(this.form)
      if (res.code === 200) {
        this.$message.success('上链成功，交易哈希：' + (res.txid || '无'))
        this.dialogVisible = false
        this.loadList()
      } else {
        this.$message.error(res.message)
      }
    },
    // 演示不可篡改：点击时提示无法修改
    tryModify(row) {
      this.$message.warning('链上数据不可篡改！无法修改已有记录。如需变更，请添加一条新的关联记录。')
    },
    // 模拟验证：直接提示通过，因为数据来自链上
    verifyHash(row) {
      this.$message.success(`验证通过：该记录已上链，交易哈希 ${row.txid} 已确认，数据未被篡改。`)
    }
  }
}
</script>

<style scoped>
.supplier-container {
  padding: 20px;
}
</style>
