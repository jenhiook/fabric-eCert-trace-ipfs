<template>
  <div class="enterprise-trace-container">
    <h3>企业端证照溯源</h3>

    <el-form :inline="true" class="search-form">
      <el-form-item label="溯源码">
        <el-input v-model="certId" placeholder="请输入溯源码" clearable style="width: 300px" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadHistory">查看历史</el-button>
        <el-button :disabled="!historyList.length" @click="exportToCSV">导出审计报表</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="pagedHistoryList" border stripe fit>
      <el-table-column prop="time" label="时间" width="180" />
      <el-table-column prop="eventType" label="事件类型" width="120" />
      <el-table-column prop="operator" label="操作人" width="120" />
      <el-table-column prop="role" label="角色" width="120" />
      <el-table-column prop="certNumber" label="证件编号" width="180" />
      <el-table-column prop="statusBefore" label="状态变更前" width="120" />
      <el-table-column prop="statusAfter" label="状态变更后" width="120" />
      <el-table-column prop="txid" label="交易哈希" min-width="200" />
    </el-table>

    <el-pagination
      background
      layout="total, sizes, prev, pager, next"
      :total="total"
      :page-size="limit"
      :current-page.sync="page"
      style="margin-top:20px; text-align:right;"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
import { getHistory } from '@/api/trace'

export default {
  data() {
    return {
      certId: '',
      historyList: [],
      total: 0,
      page: 1,
      limit: 10,
      loading: false
    }
  },
  computed: {
    pagedHistoryList() {
      const start = (this.page - 1) * this.limit
      return this.historyList.slice(start, start + this.limit)
    }
  },
  methods: {
    async loadHistory() {
      if (!this.certId) {
        this.$message.warning('请输入溯源码')
        return
      }
      this.loading = true
      try {
        const res = await getHistory(this.certId)
        let list = []
        if (res.code === 200 && res.data) {
          if (typeof res.data === 'string') {
            list = JSON.parse(res.data)
          } else {
            list = res.data
          }
        }
        this.historyList = list || []
        this.total = this.historyList.length
        this.page = 1
      } catch (e) {
        console.error(e)
        this.$message.error('加载历史失败')
      } finally {
        this.loading = false
      }
    },
    exportToCSV() {
      if (!this.historyList.length) {
        this.$message.warning('暂无数据可导出')
        return
      }
      const headers = ['时间', '事件类型', '操作人', '角色', '证件编号', '状态变更前', '状态变更后', '交易哈希']
      const rows = this.historyList.map(item => [
        item.time || '',
        item.eventType || '',
        item.operator || '',
        item.role || '',
        item.certNumber || '',
        item.statusBefore || '',
        item.statusAfter || '',
        item.txid || ''
      ])
      let csvContent = headers.join(',') + '\n'
      rows.forEach(row => {
        const escapedRow = row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(',')
        csvContent += escapedRow + '\n'
      })
      const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' })
      const link = document.createElement('a')
      const url = URL.createObjectURL(blob)
      link.href = url
      link.setAttribute('download', `audit_${this.certId}_${new Date().toISOString().slice(0, 19)}.csv`)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)
      this.$message.success('导出成功')
    },
    handleSizeChange(val) {
      this.limit = val
      this.page = 1
    },
    handleCurrentChange(val) {
      this.page = val
    }
  }
}
</script>

<style scoped>
.enterprise-trace-container {
  padding: 20px;
}
.search-form {
  margin-bottom: 20px;
}
</style>
