<template>
  <div class="dashboard-container">
    <h2 class="title">区块链电子证照溯源可视化分析平台</h2>

    <!-- 统计卡片 Flex 布局（完美对齐） -->
    <div class="stat-cards-wrapper">
      <div v-for="card in statCards" :key="card.title" class="stat-card-item">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon">
            <i :class="card.icon" />
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ card.value }}</div>
            <div class="stat-title">{{ card.title }}</div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 双图表 -->
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="chart-card">
          <div slot="header">证照阶段分布</div>
          <div id="stageChart" style="height: 300px" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <div slot="header">证照种类分布</div>
          <div id="typeChart" style="height: 300px" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 证照列表 -->
    <el-card class="table-card">
      <div slot="header">我的证照列表</div>
      <el-table :data="pagedCertList" border stripe>
        <el-table-column prop="traceability_code" label="溯源码" width="180" />
        <el-table-column prop="farmer_input.fa_fruitName" label="证照类型" />
        <el-table-column prop="farmer_input.fa_farmerName" label="持证人" />
        <el-table-column prop="farmer_input.fa_origin" label="编号" />
        <el-table-column prop="farmer_input.fa_timestamp" label="上链时间" />
        <el-table-column label="阶段状态">
          <template slot-scope="scope">{{ getStageLabel(scope.row) }}</template>
        </el-table-column>
      </el-table>
      <el-pagination
        layout="total, prev, pager, next"
        :total="certList.length"
        :page-size="pageSize"
        :current-page.sync="currentPage"
        style="margin-top:15px; text-align:right"
        @current-change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { getUserStats, getFruitList } from '@/api/trace'

export default {
  data() {
    return {
      certList: [],
      currentPage: 1,
      pageSize: 8,
      stats: {
        totalCerts: 0,
        stageDist: { 已上链: 0, 已审核: 0, 已备案: 0, 已核验: 0 },
        certTypeDist: {}
      },
      stageChart: null,
      typeChart: null
    }
  },
  computed: {
    // 使用 Element UI 最基础的图标类（一定能显示）
    statCards() {
      return [
        { title: '我的证照总数', value: this.stats.totalCerts, icon: 'el-icon-document' },
        { title: '已上链', value: this.stats.stageDist['已上链'] || 0, icon: 'el-icon-upload' },
        { title: '已审核', value: this.stats.stageDist['已审核'] || 0, icon: 'el-icon-check' },
        { title: '已备案', value: this.stats.stageDist['已备案'] || 0, icon: 'el-icon-folder-checked' },
        { title: '已核验', value: this.stats.stageDist['已核验'] || 0, icon: 'el-icon-circle-check' }
      ]
    },
    pagedCertList() {
      const start = (this.currentPage - 1) * this.pageSize
      return this.certList.slice(start, start + this.pageSize)
    }
  },
  async mounted() {
    await this.loadData()
    this.initCharts()
    this.updateCharts()
  },
  methods: {
    async loadData() {
      try {
        const res = await getUserStats()
        if (res.code === 200) this.stats = res.data
      } catch (e) { console.error(e) }
      try {
        const res = await getFruitList()
        if (res.code === 200 && res.data) {
          this.certList = JSON.parse(res.data)
          this.recalcStageDist()
        }
      } catch (e) { console.error(e) }
    },
    recalcStageDist() {
      const stageDist = { 已上链: 0, 已审核: 0, 已备案: 0, 已核验: 0 }
      const typeDist = {}
      for (const cert of this.certList) {
        const stage = this.getStageLabel(cert)
        if (stageDist[stage] !== undefined) stageDist[stage]++
        const certType = cert.farmer_input?.fa_fruitName || '未分类'
        typeDist[certType] = (typeDist[certType] || 0) + 1
      }
      this.stats.stageDist = stageDist
      this.stats.certTypeDist = typeDist
      this.stats.totalCerts = this.certList.length
      this.updateCharts()
    },
    getStageLabel(row) {
      const s = row.shop_input || {}
      const d = row.driver_input || {}
      const f = row.factory_input || {}
      const p = row.farmer_input || {}
      if (s.sh_txid) return '已核验'
      if (d.dr_txid) return '已备案'
      if (f.fac_txid) return f.fac_contactNumber === '通过' ? '已审核' : '审核驳回'
      if (p.fa_txid) return '已上链'
      return '--'
    },
    initCharts() {
      this.stageChart = echarts.init(document.getElementById('stageChart'))
      this.typeChart = echarts.init(document.getElementById('typeChart'))
    },
    updateCharts() {
      if (!this.stageChart || !this.typeChart) return
      const stageOrder = ['已上链', '已审核', '已备案', '已核验']
      const stageData = stageOrder.map(name => ({ name, value: this.stats.stageDist[name] || 0 }))

      this.stageChart.setOption({
        tooltip: { trigger: 'item' },
        legend: { top: 'bottom' },
        series: [{ type: 'pie', radius: '55%', data: stageData, emphasis: { scale: true }}]
      })

      this.typeChart.setOption({
        tooltip: { trigger: 'axis' },
        xAxis: { type: 'category', data: Object.keys(this.stats.certTypeDist), axisLabel: { rotate: 30 }},
        yAxis: { type: 'value' },
        series: [{ type: 'bar', data: Object.values(this.stats.certTypeDist), itemStyle: { borderRadius: [4, 4, 0, 0], color: '#409eff' }}]
      })
    },
    handlePageChange(page) {
      this.currentPage = page
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 20px;
  background: #f0f2f5;
  min-height: 100vh;
}
.title {
  text-align: center;
  margin-bottom: 24px;
  font-weight: 500;
  color: #303133;
}
/* Flex 布局：5个卡片等宽均匀排列 */
.stat-cards-wrapper {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  margin-bottom: 20px;
}
.stat-card-item {
  flex: 1;
  margin: 0 8px;
  min-width: 150px;
}
.stat-card-item:first-child {
  margin-left: 0;
}
.stat-card-item:last-child {
  margin-right: 0;
}
.stat-card {
  text-align: center;
  border-radius: 16px;
  transition: all 0.3s;
}
.stat-icon {
  font-size: 40px;
  color: #409eff;
  margin-bottom: 8px;
}
.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}
.stat-title {
  font-size: 14px;
  color: #909399;
  margin-top: 6px;
}
.chart-card, .table-card {
  border-radius: 16px;
  margin-bottom: 20px;
  overflow: hidden;
}
@media (max-width: 768px) {
  .stat-card-item {
    flex: 0 0 calc(50% - 16px);
    margin: 8px;
  }
  .stat-card-item:first-child, .stat-card-item:last-child {
    margin: 8px;
  }
}
</style>
