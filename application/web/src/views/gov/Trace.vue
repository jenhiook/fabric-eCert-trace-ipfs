<template>
  <div class="gov-trace-container">
    <h3>政府端证照溯源</h3>

    <el-form :inline="true" class="search-form">
      <el-form-item label="溯源码">
        <el-input v-model="certId" placeholder="请输入溯源码" clearable style="width: 240px" />
      </el-form-item>
      <el-form-item label="起始时间">
        <el-date-picker v-model="startTime" type="datetime" placeholder="选择起始时间" value-format="timestamp" clearable />
      </el-form-item>
      <el-form-item label="结束时间">
        <el-date-picker v-model="endTime" type="datetime" placeholder="选择结束时间" value-format="timestamp" clearable />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="searchHistory">查询</el-button>
        <el-button @click="exportToCSV">导出审计报表</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="historyList" border stripe>
      <el-table-column prop="time" label="时间" width="180" />
      <el-table-column prop="eventType" label="事件类型" width="120" />
      <el-table-column prop="operator" label="操作人" width="120" />
      <el-table-column prop="role" label="角色" width="120" />
      <el-table-column prop="certNumber" label="证件编号" width="180" />
      <el-table-column prop="statusBefore" label="状态变更前" width="120" />
      <el-table-column prop="statusAfter" label="状态变更后" width="120" />
      <el-table-column prop="txid" label="交易哈希" width="200" />
    </el-table>

    <el-pagination
      background
      layout="total, sizes, prev, pager, next"
      :total="total"
      :page-size="limit"
      :current-page="page"
      style="margin-top:20px; text-align:right;"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
import { govGetHistory } from '@/api/gov'

export default {
  data() {
    return {
      certId: '',
      startTime: null,
      endTime: null,
      historyList: [],
      total: 0,
      page: 1,
      limit: 10
    }
  },
  created() {
    const code = this.$route.query.code
    if (code) {
      this.certId = code
      this.searchHistory()
    }
  },
  methods: {
    async searchHistory() {
      if (!this.certId) {
        this.$message.warning('请输入溯源码')
        return
      }
      const params = {
        certId: this.certId,
        page: this.page,
        limit: this.limit,
        startTime: this.startTime ? Math.floor(this.startTime / 1000) : '',
        endTime: this.endTime ? Math.floor(this.endTime / 1000) : ''
      }
      try {
        const res = await govGetHistory(params)
        if (res.code === 200) {
          let historyArray = []
          let totalCount = 0
          if (Array.isArray(res.data)) {
            historyArray = res.data
            totalCount = res.data.length
          } else if (res.data && Array.isArray(res.data.history)) {
            historyArray = res.data.history
            totalCount = res.data.total || historyArray.length
          }
          this.historyList = historyArray
          this.total = totalCount
        } else {
          this.$message.error(res.message || '查询失败')
        }
      } catch (e) {
        console.error(e)
        this.$message.error('请求异常')
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
      this.searchHistory()
    },
    handleCurrentChange(val) {
      this.page = val
      this.searchHistory()
    }
  }
}
</script>

<style scoped>
.gov-trace-container {
  padding: 20px;
}
.search-form {
  margin-bottom: 20px;
}
</style>
