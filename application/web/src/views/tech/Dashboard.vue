<template>
  <div class="dashboard-container">
    <h2 class="title">技术实体端证照统计分析</h2>

    <!-- 个人核验总数卡片 -->
    <div class="stat-cards-wrapper">
      <div class="stat-card-item">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon"><i class="el-icon-document" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ myCertList.length }}</div>
            <div class="stat-title">我核验的证照总数</div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 全局阶段分布卡片 -->
    <div class="stat-cards-wrapper">
      <div
        v-for="card in stageCards"
        :key="card.title"
        class="stat-card-item"
      >
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon"><i :class="card.icon" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ card.value }}</div>
            <div class="stat-title">{{ card.title }}</div>
          </div>
        </el-card>
      </div>
    </div>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="chart-card">
          <div slot="header">全局证照阶段分布</div>
          <div id="stageChart" style="height: 300px" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <div slot="header">全局证照种类分布</div>
          <div id="typeChart" style="height: 300px" />
        </el-card>
      </el-col>
    </el-row>

    <el-card class="table-card">
      <div slot="header">我核验的证照列表</div>
      <el-table :data="pagedCertList" border stripe>
        <el-table-column prop="traceability_code" label="溯源码" min-width="180" />
        <el-table-column prop="farmer_input.fa_fruitName" label="证照类型" min-width="110" />
        <el-table-column prop="farmer_input.fa_farmerName" label="持证人" min-width="100" />
        <el-table-column prop="farmer_input.fa_origin" label="编号" min-width="180" />
        <el-table-column prop="shop_input.sh_timestamp" label="核验时间" min-width="160" />
        <el-table-column label="核验结果" min-width="90">
          <template slot-scope="scope">
            <el-tag :type="scope.row.shop_input.sh_shopAddress === '有效' ? 'success' : 'danger'">
              {{ scope.row.shop_input.sh_shopAddress === '有效' ? '有效' : '无效' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="100">
          <template slot-scope="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        layout="total, prev, pager, next"
        :total="myCertList.length"
        :page-size="pageSize"
        :current-page.sync="currentPage"
        style="margin-top:15px; text-align:right;"
        @current-change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { getAllFruitInfo } from '@/api/trace'

export default {
  data() {
    return {
      allCerts: [],
      myCertList: [],
      currentPage: 1,
      pageSize: 8,
      stageDist: { 已上链: 0, 已审核: 0, 已备案: 0, 已核验: 0 },
      certTypeDist: {},
      stageChart: null,
      typeChart: null
    }
  },
  computed: {
    stageCards() {
      return [
        { title: '已上链', value: this.stageDist['已上链'] || 0, icon: 'el-icon-upload' },
        { title: '已审核', value: this.stageDist['已审核'] || 0, icon: 'el-icon-check' },
        { title: '已备案', value: this.stageDist['已备案'] || 0, icon: 'el-icon-folder-checked' },
        { title: '已核验', value: this.stageDist['已核验'] || 0, icon: 'el-icon-circle-check' }
      ]
    },
    pagedCertList() {
      const start = (this.currentPage - 1) * this.pageSize
      return this.myCertList.slice(start, start + this.pageSize)
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
        const res = await getAllFruitInfo()
        if (res.code === 200 && res.data) {
          this.allCerts = JSON.parse(res.data)
          // 个人核验列表：shop_input.sh_txid 存在且不为空
          this.myCertList = this.allCerts.filter(cert => {
            const shop = cert.shop_input || {}
            return shop.sh_txid && shop.sh_txid !== ''
          })
          // 全局统计
          const stage = { 已上链: 0, 已审核: 0, 已备案: 0, 已核验: 0 }
          const type = {}
          for (const cert of this.allCerts) {
            // 阶段判定
            const s = this.getStageLabel(cert)
            if (stage[s] !== undefined) stage[s]++
            // 种类判定
            const t = (cert.farmer_input && cert.farmer_input.fa_fruitName) || '未分类'
            type[t] = (type[t] || 0) + 1
          }
          this.stageDist = stage
          this.certTypeDist = type
        }
      } catch (e) {
        console.error(e)
        this.$message.error('加载数据失败')
      }
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
      const stageData = stageOrder.map(name => ({ name, value: this.stageDist[name] || 0 }))

      this.stageChart.setOption({
        tooltip: { trigger: 'item' },
        legend: { top: 'bottom' },
        series: [{ type: 'pie', radius: '55%', data: stageData, emphasis: { scale: true }}]
      })

      this.typeChart.setOption({
        tooltip: { trigger: 'axis' },
        xAxis: { type: 'category', data: Object.keys(this.certTypeDist), axisLabel: { rotate: 30 }},
        yAxis: { type: 'value' },
        series: [{ type: 'bar', data: Object.values(this.certTypeDist), itemStyle: { borderRadius: [4, 4, 0, 0], color: '#409eff' }}]
      })
    },
    handlePageChange(page) {
      this.currentPage = page
    },
    viewDetail(row) {
      this.$router.push({ path: '/trace', query: { code: row.traceability_code }})
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
.stat-cards-wrapper {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 16px;
  margin-bottom: 20px;
}
.stat-card-item {
  flex: 1;
  min-width: 150px;
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
  }
}
</style>
