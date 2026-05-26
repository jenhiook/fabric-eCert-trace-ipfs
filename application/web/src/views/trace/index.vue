<template>
  <div class="trace-container">
    <h3>电子证照溯源查询</h3>

    <div class="search-bar">
      <el-input v-model="traceCode" placeholder="请输入溯源码" style="width:340px" clearable @keyup.enter.native="queryInfo" />
      <el-button type="primary" icon="el-icon-search" style="margin-left:12px" @click="queryInfo">溯源查询</el-button>
      <el-button icon="el-icon-refresh" @click="getAll">获取所有证照</el-button>
    </div>

    <div v-if="cert">
      <el-button icon="el-icon-back" size="small" @click="goBack">返回列表</el-button>
      <!-- 独立时间轴按钮：仅非政务用户显示 -->
      <el-button
        v-if="$store.state.user.userType == '个人用户'"
        icon="el-icon-time"
        size="small"
        @click="gotoTimeline"
      >独立时间轴</el-button>

      <div class="cert-header">
        <h4>证照全流程追踪</h4>
        <p>溯源码：<span class="code">{{ cert.traceability_code }}</span></p>
      </div>

      <div class="cert-preview" @click="openImg(certImgSrc)">
        <img :src="certImgSrc" alt="证照" class="preview-img">
        <p class="img-hint">点击查看大图</p>
      </div>

      <!-- 四阶段卡片（保持不变） -->
      <el-timeline>
        <el-timeline-item placement="top" type="primary">
          <el-card shadow="hover" class="stage-card" :class="{ 'stage-active': hasField('personal') }">
            <div class="stage-header">
              <span class="stage-tag tag-chain">1</span>
              <span class="stage-title">证照上链 · 个人用户录入</span>
              <el-tag v-if="hasField('personal')" type="success" size="small">已完成</el-tag>
              <el-tag v-else type="info" size="small">未录入</el-tag>
            </div>
            <el-divider />
            <div v-if="hasField('personal')" class="info-grid">
              <div class="info-item"><span class="info-label">证照类型</span><span class="info-value">{{ pi.fa_fruitName }}</span></div>
              <div class="info-item"><span class="info-label">持证人</span><span class="info-value">{{ pi.fa_farmerName }}</span></div>
              <div class="info-item"><span class="info-label">性别</span><span class="info-value">{{ pi.fa_plantTime }}</span></div>
              <div class="info-item"><span class="info-label">证照编号</span><span class="info-value">{{ pi.fa_origin }}</span></div>
              <div class="info-item"><span class="info-label">联系电话</span><span class="info-value">{{ pi.fa_pickingTime }}</span></div>
              <div class="info-item"><span class="info-label">上链时间</span><span class="info-value">{{ pi.fa_timestamp }}</span></div>
              <div v-if="pi.fa_certImage && pi.fa_certImage.startsWith('ipfs://')" class="info-item full-width ipfs-row">
                <span class="info-label">IPFS 分布式存储标识（CID）</span>
                <span class="info-value ipfs-cid">
                  {{ pi.fa_certImage.replace('ipfs://', '') }}
                  <el-button type="text" size="mini" icon="el-icon-document-copy" @click="copyCID(pi.fa_certImage)">复制 CID</el-button>
                </span>
              </div>
              <div class="info-item full-width"><span class="info-label">交易哈希</span><span class="info-value tx-hash">{{ pi.fa_txid }}</span></div>
            </div>
            <p v-else class="stage-empty">尚未录入</p>
          </el-card>
        </el-timeline-item>

        <el-timeline-item placement="top" type="warning">
          <el-card shadow="hover" class="stage-card" :class="{ 'stage-active': hasField('govt') }">
            <div class="stage-header">
              <span class="stage-tag tag-audit">2</span>
              <span class="stage-title">政务审核 · 部门审批</span>
              <el-tag v-if="hasField('govt')" :type="fi2.fac_contactNumber === '通过' ? 'success' : 'danger'" size="small">
                {{ fi2.fac_contactNumber === '通过' ? '审核通过' : '审核驳回' }}
              </el-tag>
              <el-tag v-else type="info" size="small">待审核</el-tag>
            </div>
            <el-divider />
            <div v-if="hasField('govt')" class="info-grid">
              <div class="info-item"><span class="info-label">部门名称</span><span class="info-value">{{ fi2.fac_productName }}</span></div>
              <div class="info-item"><span class="info-label">部门代码</span><span class="info-value">{{ fi2.fac_productionbatch }}</span></div>
              <div class="info-item"><span class="info-label">审核人</span><span class="info-value">{{ fi2.fac_factoryName }}</span></div>
              <div class="info-item"><span class="info-label">审核结果</span><span class="info-value">{{ fi2.fac_contactNumber }}</span></div>
              <div class="info-item"><span class="info-label">审核时间</span><span class="info-value">{{ fi2.fac_timestamp }}</span></div>
              <div class="info-item full-width"><span class="info-label">交易哈希</span><span class="info-value tx-hash">{{ fi2.fac_txid }}</span></div>
            </div>
            <p v-else class="stage-empty">尚未审核</p>
          </el-card>
        </el-timeline-item>

        <el-timeline-item placement="top" type="info">
          <el-card shadow="hover" class="stage-card" :class="{ 'stage-active': hasField('enterprise') }">
            <div class="stage-header">
              <span class="stage-tag tag-file">3</span>
              <span class="stage-title">企业备案 · 使用记录</span>
              <el-tag v-if="hasField('enterprise')" type="success" size="small">已备案</el-tag>
              <el-tag v-else-if="!hasField('govt') || fi2.fac_contactNumber !== '通过'" type="warning" size="small">前置未完成</el-tag>
              <el-tag v-else type="info" size="small">未备案</el-tag>
            </div>
            <el-divider />
            <div v-if="hasField('enterprise')" class="info-grid">
              <div class="info-item"><span class="info-label">企业名称</span><span class="info-value">{{ di.dr_name }}</span></div>
              <div class="info-item"><span class="info-label">信用代码</span><span class="info-value">{{ di.dr_age }}</span></div>
              <div class="info-item"><span class="info-label">使用目的</span><span class="info-value">{{ di.dr_phone }}</span></div>
              <div class="info-item"><span class="info-label">经办人</span><span class="info-value">{{ di.dr_transport }}</span></div>
              <div class="info-item"><span class="info-label">备案时间</span><span class="info-value">{{ di.dr_timestamp }}</span></div>
              <div class="info-item full-width"><span class="info-label">交易哈希</span><span class="info-value tx-hash">{{ di.dr_txid }}</span></div>
            </div>
            <p v-else class="stage-empty">{{ !hasField('govt') || fi2.fac_contactNumber !== '通过' ? '请先通过政务审核' : '尚未备案' }}</p>
          </el-card>
        </el-timeline-item>

        <el-timeline-item placement="top" type="success">
          <el-card shadow="hover" class="stage-card" :class="{ 'stage-active': hasField('tech') }">
            <div class="stage-header">
              <span class="stage-tag tag-verify">4</span>
              <span class="stage-title">技术核验 · 安全验证</span>
              <el-tag v-if="hasField('tech')" :type="si.sh_shopAddress === '有效' ? 'success' : 'danger'" size="small">
                {{ si.sh_shopAddress === '有效' ? '核验有效' : '核验无效' }}
              </el-tag>
              <el-tag v-else-if="!hasField('enterprise')" type="warning" size="small">前置未完成</el-tag>
              <el-tag v-else type="info" size="small">未核验</el-tag>
            </div>
            <el-divider />
            <div v-if="hasField('tech')" class="info-grid">
              <div class="info-item"><span class="info-label">核验实体</span><span class="info-value">{{ si.sh_storeTime }}</span></div>
              <div class="info-item"><span class="info-label">服务类型</span><span class="info-value">{{ si.sh_sellTime }}</span></div>
              <div class="info-item"><span class="info-label">安全等级</span><span class="info-value">{{ si.sh_shopName }}</span></div>
              <div class="info-item"><span class="info-label">核验结果</span><span class="info-value">{{ si.sh_shopAddress }}</span></div>
              <div class="info-item"><span class="info-label">核验时间</span><span class="info-value">{{ si.sh_timestamp }}</span></div>
              <div class="info-item full-width"><span class="info-label">交易哈希</span><span class="info-value tx-hash">{{ si.sh_txid }}</span></div>
            </div>
            <p v-else class="stage-empty">{{ !hasField('enterprise') ? '请先完成企业备案' : '尚未核验' }}</p>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </div>

    <!-- 证照列表（保持不变） -->
    <el-card v-if="allList.length > 0 && !cert" style="margin-top:20px">
      <div slot="header">全部证照列表（共 {{ allList.length }} 条）</div>
      <el-table :data="pagedList" border stripe style="cursor:pointer" @row-click="rowClick">
        <el-table-column prop="traceability_code" label="溯源码" width="180" />
        <el-table-column prop="farmer_input.fa_fruitName" label="证照类型" width="110" />
        <el-table-column prop="farmer_input.fa_origin" label="编号" width="200" />
        <el-table-column prop="farmer_input.fa_farmerName" label="姓名" width="100" />
        <el-table-column prop="farmer_input.fa_timestamp" label="上链时间" width="180" />
        <el-table-column label="当前阶段" width="120">
          <template slot-scope="scope">{{ stageLabel(scope.row) }}</template>
        </el-table-column>
        <el-table-column label="证照预览" width="90">
          <template slot-scope="scope">
            <el-button v-if="hasImg(scope.row)" type="text" size="small" @click.stop="openImg(getDisplayImg((scope.row.farmer_input || {}).fa_certImage))">查看</el-button>
            <span v-else>无</span>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        style="margin-top:15px;text-align:right"
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[5, 10, 20]"
        layout="total, sizes, prev, pager, next"
        :total="allList.length"
        @size-change="s => { pageSize = s; page = 1 }"
        @current-change="p => page = p"
      />
    </el-card>

    <el-dialog :visible.sync="showImg" width="460px" :show-header="false" :close-on-click-modal="true">
      <img :src="viewImg" style="width:100%;border-radius:6px">
    </el-dialog>
  </div>
</template>

<script>
import { getFruitInfo, getAllFruitInfo, getFruitList, getHistory } from '@/api/trace'
const defaultImg = require('@/assets/camera.jpg')

export default {
  data() {
    return {
      traceCode: '',
      cert: null,
      allList: [],
      showImg: false,
      viewImg: '',
      page: 1,
      pageSize: 10,
      historyList: []
    }
  },
  computed: {
    certImgSrc() {
      const img = (this.cert && this.cert.farmer_input && this.cert.farmer_input.fa_certImage) || ''
      if (img.startsWith('data:')) return img
      if (img.startsWith('ipfs://')) return 'http://127.0.0.1:8080/ipfs/' + img.replace('ipfs://', '')
      return defaultImg
    },
    pi() { return (this.cert && this.cert.farmer_input) || {} },
    fi2() { return (this.cert && this.cert.factory_input) || {} },
    di() { return (this.cert && this.cert.driver_input) || {} },
    si() { return (this.cert && this.cert.shop_input) || {} },
    pagedList() {
      const start = (this.page - 1) * this.pageSize
      return this.allList.slice(start, start + this.pageSize)
    }
  },
  // 监听路由参数自动查询（支持从证照审核页面跳转）
  watch: {
    '$route.query.code': {
      handler(newCode) {
        if (newCode) {
          this.traceCode = newCode
          this.queryInfo()
        }
      },
      immediate: true
    }
  },
  methods: {
    hasField(type) {
      const map = { personal: 'fa_txid', govt: 'fac_txid', enterprise: 'dr_txid', tech: 'sh_txid' }
      const data = { personal: this.pi, govt: this.fi2, enterprise: this.di, tech: this.si }
      return !!(data[type] && data[type][map[type]])
    },
    hasImg(row) {
      const img = (row.farmer_input || {}).fa_certImage || ''
      return img.startsWith('data:') || img.startsWith('ipfs://')
    },
    getDisplayImg(img) {
      if (!img) return ''
      if (img.startsWith('ipfs://')) return 'http://127.0.0.1:8080/ipfs/' + img.replace('ipfs://', '')
      return img
    },
    stageLabel(row) {
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
    queryInfo() {
      if (!this.traceCode) return this.$message.error('请输入溯源码')
      const fd = new FormData()
      fd.append('traceability_code', this.traceCode)
      getFruitInfo(fd).then(res => {
        if (res.code === 200) {
          this.cert = JSON.parse(res.data)
          this.allList = []
          this.page = 1
          this.loadHistory(this.cert.traceability_code)
        } else {
          this.$message.error(res.message || '查询失败')
        }
      })
    },
    gotoTimeline() {
      if (this.cert && this.cert.traceability_code) {
        this.$router.push({
          path: '/timeline',
          query: { code: this.cert.traceability_code }
        })
      }
    },
    getAll() {
      const user = this.$store.state.user
      if (user && user.userType === '个人用户') {
        getFruitList().then(res => {
          if (res.code === 200 && res.data) {
            const list = JSON.parse(res.data)
            list.sort((a, b) => {
              return ((b.farmer_input || {}).fa_timestamp || '').localeCompare((a.farmer_input || {}).fa_timestamp || '')
            })
            this.allList = list
            this.cert = null
            this.page = 1
          }
        })
      } else {
        getAllFruitInfo().then(res => {
          if (res.code === 200) {
            const list = JSON.parse(res.data)
            list.sort((a, b) => {
              return ((b.farmer_input || {}).fa_timestamp || '').localeCompare((a.farmer_input || {}).fa_timestamp || '')
            })
            this.allList = list
            this.cert = null
            this.page = 1
          }
        })
      }
    },
    goBack() {
      this.cert = null
      if (this.allList.length === 0) this.getAll()
    },
    rowClick(row) {
      this.traceCode = row.traceability_code
      this.cert = row
      this.allList = []
      this.loadHistory(row.traceability_code)
    },
    copyCID(url) {
      const cid = url.replace('ipfs://', '')
      navigator.clipboard.writeText(cid).then(() => {
        this.$message.success('IPFS CID 已复制到剪贴板')
      })
    },
    copyText(text) {
      if (!text) return
      navigator.clipboard.writeText(text).then(() => {
        this.$message.success('已复制到剪贴板')
      })
    },
    openImg(url) {
      if (!url) return
      this.viewImg = url
      this.showImg = true
    },
    async loadHistory(certId) {
      try {
        const res = await getHistory(certId)
        console.log('📥 后端返回历史数据:', res)

        let list = []
        if (res.code === 200 && res.data) {
          if (typeof res.data === 'string') {
            try {
              list = JSON.parse(res.data)
            } catch (e) {
              list = []
            }
          } else {
            list = res.data
          }
        }

        if (!Array.isArray(list)) {
          list = []
        }

        this.historyList = list.map(item => ({
          eventType: item.eventType || '',
          operator: item.operator || '',
          role: item.role || '',
          certNumber: item.certNumber || '',
          txid: item.txid || '',
          statusBefore: item.statusBefore || '',
          statusAfter: item.statusAfter || '',
          cid: item.cid || '',
          time: this.formatTime(item.time)
        })).filter(item => item.eventType)
          .sort((a, b) => new Date(a.time) - new Date(b.time))
      } catch (e) {
        console.error('❌ loadHistory异常:', e)
        this.historyList = []
      }
    },
    formatTime(timestamp) {
      if (!timestamp) return ''
      let ts = parseInt(timestamp)
      if (isNaN(ts)) return timestamp
      if (ts < 10000000000) ts *= 1000
      const date = new Date(ts)
      return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour12: false })}`
    },
    shortenTxid(txid) {
      if (!txid) return ''
      if (txid.length <= 20) return txid
      return txid.slice(0, 10) + '...' + txid.slice(-10)
    },
    timelineType(eventType) {
      const map = { '证照上链': 'primary', '政务审核': 'warning', '企业备案': 'info', '技术核验': 'success' }
      return map[eventType] || 'primary'
    },
    timelineIcon(eventType) {
      const map = { '证照上链': 'el-icon-upload', '政务审核': 'el-icon-check', '企业备案': 'el-icon-document', '技术核验': 'el-icon-shield' }
      return map[eventType] || 'el-icon-info'
    },
    eventTagType(eventType) {
      const map = { '证照上链': 'primary', '政务审核': 'warning', '企业备案': 'info', '技术核验': 'success' }
      return map[eventType] || ''
    },
    statusAfterType(status) {
      if (!status) return 'info'
      if (status.includes('通过') || status.includes('有效') || status.includes('备案') || status.includes('核验')) return 'success'
      if (status.includes('驳回') || status.includes('无效')) return 'danger'
      return 'info'
    }
  }
}
</script>

<style scoped>
/* 原有样式保留 */
.trace-container { padding: 25px 40px; max-width: 1100px; margin: 0 auto; }
h3 { text-align: center; margin-bottom: 12px; }
.search-bar { text-align: center; margin-bottom: 20px; }
.cert-header { text-align: center; margin: 15px 0 10px; }
.cert-header h4 { margin: 0 0 4px; }
.code { color: #409eff; font-family: monospace; font-size: 13px; }
.cert-preview { text-align: center; margin-bottom: 18px; cursor: pointer; }
.preview-img { width: 90px; height: 120px; object-fit: cover; border-radius: 6px; border: 2px dashed #ddd; }
.img-hint { font-size: 11px; color: #bbb; margin-top: 3px; }
.stage-card { margin-bottom: 6px; }
.stage-active { border-left: 3px solid #409eff; }
.stage-header { display: flex; align-items: center; gap: 8px; margin-bottom: 0; }
.stage-tag { display: inline-block; width: 22px; height: 22px; line-height: 22px; text-align: center; border-radius: 50%; font-size: 11px; color: #fff; font-weight: bold; }
.tag-chain { background: #409eff; }
.tag-audit { background: #e6a23c; }
.tag-file { background: #909399; }
.tag-verify { background: #67c23a; }
.stage-title { font-size: 14px; font-weight: 600; }
.stage-empty { color: #bbb; font-style: italic; padding: 8px 0; }
.info-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 0; border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden; }
.info-item { display: flex; flex-direction: column; padding: 8px 12px; min-height: 52px; border-right: 1px solid #ebeef5; border-bottom: 1px solid #ebeef5; }
.info-item.full-width { grid-column: 1 / -1; }
.info-label { font-size: 11px; color: #999; margin-bottom: 3px; }
.info-value { font-size: 13px; color: #333; word-break: break-all; }
.tx-hash { font-family: monospace; font-size: 11px; color: #303133; }
.ipfs-row { background: #f0f9f5; }
.ipfs-cid { font-family: monospace; font-size: 12px; color: #5a9e9b; display: flex; align-items: center; gap: 8px; }

/* 原有时间轴卡片样式已删除，只保留必要样式 */
.history-card { margin-bottom: 8px; }
.history-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; flex-wrap: wrap; gap: 8px; }
.history-operator { font-size: 14px; color: #606266; }
.role-badge { color: #909399; font-size: 12px; }
.history-details { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.detail-item { display: flex; flex-direction: column; }
.detail-item.full-width { grid-column: 1 / -1; }
.detail-label { font-size: 12px; color: #909399; margin-bottom: 4px; }
.detail-value { font-size: 13px; color: #303133; word-break: break-all; }
.status-flow { flex-direction: row; align-items: center; justify-content: space-between; }
.status-badge { display: flex; align-items: center; gap: 8px; }
.el-icon-right { color: #c0c4cc; }
</style>
