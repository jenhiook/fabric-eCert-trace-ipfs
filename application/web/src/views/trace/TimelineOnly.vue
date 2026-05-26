<template>
  <div class="timeline-only-container">
    <h3>证照历史时间轴（区块链全记录）</h3>

    <div class="search-bar">
      <el-input
        v-model="traceCode"
        placeholder="请输入溯源码"
        style="width:340px"
        clearable
        @keyup.enter.native="searchAndLoad"
      />
      <el-button type="primary" icon="el-icon-search" @click="searchAndLoad">查询</el-button>
      <el-button icon="el-icon-download" :disabled="historyList.length === 0" @click="exportAuditLog">导出审计日志</el-button>
    </div>

    <div v-if="certInfo" class="cert-summary">
      <el-card>
        <div slot="header">证照摘要</div>
        <div class="summary-row">
          <span>持证人：{{ certInfo.holderName }}</span>
          <span>证照编号：{{ certInfo.certNumber }}</span>
          <el-button type="text" icon="el-icon-picture" @click="openCertImage">查看原图</el-button>
        </div>
        <div v-if="certInfo.cid" class="summary-row">
          <span>IPFS CID：{{ certInfo.cid }}</span>
          <el-button type="text" icon="el-icon-document-copy" @click="copyText(certInfo.cid)">复制</el-button>
          <el-button type="text" icon="el-icon-search" @click="verifyImage">核验图片指纹</el-button>
        </div>
      </el-card>
    </div>

    <el-card v-if="historyList.length > 0" class="timeline-card">
      <div slot="header">
        <span>全流程溯源时间轴（区块链历史记录）</span>
        <span style="float:right; font-size:12px;">共 {{ historyList.length }} 条记录</span>
      </div>
      <el-timeline>
        <el-timeline-item
          v-for="(item, idx) in historyList"
          :key="item.txid || idx"
          :timestamp="item.time"
          placement="top"
          :type="timelineType(item.eventType)"
          :icon="timelineIcon(item.eventType)"
        >
          <el-card shadow="hover" class="history-card">
            <div class="history-header">
              <el-tag :type="eventTagType(item.eventType)" effect="dark" size="medium">{{ item.eventType }}</el-tag>
              <span class="history-operator"><i class="el-icon-user" /> {{ item.operator }}（{{ item.role }}）</span>
            </div>
            <div class="history-details">
              <div class="detail-item"><span class="detail-label">证件编号</span><span class="detail-value">{{ item.certNumber || '—' }}</span></div>
              <div class="detail-item">
                <span class="detail-label">事件ID（TxID）</span>
                <span class="detail-value tx-hash">{{ shortenTxid(item.txid) }}<el-button v-if="item.txid" type="text" size="mini" icon="el-icon-document-copy" @click="copyText(item.txid)">复制</el-button></span>
              </div>
              <div class="detail-item full-width">
                <span class="detail-label">状态变更</span>
                <div class="status-badge"><el-tag size="small" type="info">{{ item.statusBefore || '—' }}</el-tag><i class="el-icon-right" /><el-tag size="small" :type="statusAfterType(item.statusAfter)">{{ item.statusAfter || '—' }}</el-tag></div>
              </div>
              <div v-if="item.cid" class="detail-item full-width">
                <span class="detail-label">IPFS CID</span><span class="detail-value">{{ item.cid }} <el-button type="text" icon="el-icon-document-copy" @click="copyText(item.cid)">复制</el-button></span>
              </div>
            </div>
            <div class="evid-btn">
              <el-button size="small" type="text" icon="el-icon-document" @click="showEvidence(item)">查看证据</el-button>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <el-card v-else class="empty-card">
      <div class="empty-placeholder"><i class="el-icon-info" /><p>暂无历史记录</p><p style="font-size:12px;">请输入一个已完成全流程的证照溯源码</p></div>
    </el-card>

    <el-dialog title="证据详情" :visible.sync="evidenceDialogVisible" width="500px">
      <div v-if="evidenceItem">
        <p><strong>事件类型：</strong>{{ evidenceItem.eventType }}</p>
        <p><strong>交易哈希：</strong>{{ evidenceItem.txid }}</p>
        <p><strong>IPFS CID：</strong>{{ evidenceItem.cid || (certInfo ? certInfo.cid : '') }}</p>
        <p><strong>核验说明：</strong>下载原图并计算 SHA-256 与链上指纹比对（需链码存储指纹）。</p>
        <el-button type="primary" size="small" @click="openCertImage">查看原图</el-button>
        <el-button size="small" @click="copyText(evidenceItem.cid || (certInfo ? certInfo.cid : ''))">复制 CID</el-button>
      </div>
    </el-dialog>

    <el-dialog :visible.sync="imgDialogVisible" width="460px" :show-header="false">
      <img :src="imgPreviewSrc" style="width:100%; border-radius:6px;">
    </el-dialog>
  </div>
</template>

<script>
import { getHistory, getFruitInfo } from '@/api/trace'

export default {
  data() {
    return {
      traceCode: '',
      historyList: [],
      certInfo: null,
      evidenceDialogVisible: false,
      evidenceItem: null,
      imgDialogVisible: false,
      imgPreviewSrc: ''
    }
  },
  created() {
    const code = this.$route.query.code
    if (code) {
      this.traceCode = code
      this.searchAndLoad()
    }
  },
  methods: {
    async searchAndLoad() {
      if (!this.traceCode) {
        this.$message.warning('请输入溯源码')
        return
      }
      await this.loadCertInfo()
      await this.loadHistory()
    },
    async loadCertInfo() {
      try {
        const fd = new FormData()
        fd.append('traceability_code', this.traceCode)
        const res = await getFruitInfo(fd)
        if (res.code === 200) {
          const cert = JSON.parse(res.data)
          const farmer = cert.farmer_input || {}
          this.certInfo = {
            holderName: farmer.fa_farmerName || '',
            certNumber: farmer.fa_origin || '',
            cid: farmer.fa_certImage ? farmer.fa_certImage.replace('ipfs://', '') : '',
            imgSrc: farmer.fa_certImage || ''
          }
        } else {
          this.certInfo = null
        }
      } catch (e) {
        console.error('加载证照信息失败', e)
        this.certInfo = null
      }
    },
    async loadHistory() {
      try {
        const res = await getHistory(this.traceCode)
        console.log('历史接口返回:', res)
        let list = []
        if (res.code === 200 && res.data) {
          if (typeof res.data === 'string') {
            try { list = JSON.parse(res.data) } catch (e) { list = [] }
          } else {
            list = res.data
          }
        }
        this.historyList = list
          .filter(item => item.eventType)
          .map(item => ({
            eventType: item.eventType || '',
            operator: item.operator || '',
            role: item.role || '',
            certNumber: item.certNumber || '',
            txid: item.txid || '',
            statusBefore: item.statusBefore || '',
            statusAfter: item.statusAfter || '',
            cid: item.cid || '',
            fingerprint: item.fingerprint || '',
            note: item.note || '',
            time: this.formatTime(item.time)
          }))
          .sort((a, b) => new Date(a.time) - new Date(b.time))

        if (this.historyList.length === 0) this.$message.info('该证照暂无历史记录')
      } catch (e) {
        console.error(e)
        this.$message.error('查询失败')
        this.historyList = []
      }
    },
    formatTime(timestamp) {
      if (!timestamp) return ''
      let ts = parseInt(timestamp)
      if (isNaN(ts)) return timestamp
      if (ts < 10000000000) ts *= 1000
      return new Date(ts).toLocaleString()
    },
    shortenTxid(txid) {
      if (!txid) return ''
      return txid.length > 20 ? txid.slice(0, 10) + '...' + txid.slice(-10) : txid
    },
    copyText(text) {
      if (!text) return
      navigator.clipboard.writeText(text).then(() => this.$message.success('已复制'))
    },
    exportAuditLog() {
      if (this.historyList.length === 0) return
      const logs = this.historyList.map(item => ({
        时间: item.time,
        事件类型: item.eventType,
        操作人: item.operator,
        角色: item.role,
        证件编号: item.certNumber,
        交易哈希: item.txid,
        状态变更前: item.statusBefore,
        状态变更后: item.statusAfter,
        CID: item.cid || ''
      }))
      const csv = this.convertToCSV(logs)
      const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
      const link = document.createElement('a')
      const url = URL.createObjectURL(blob)
      link.href = url
      link.setAttribute('download', `audit_${this.traceCode}.csv`)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)
    },
    convertToCSV(arr) {
      if (!arr.length) return ''
      const header = Object.keys(arr[0]).join(',')
      const rows = arr.map(obj => Object.values(obj).map(v => `"${v}"`).join(','))
      return [header, ...rows].join('\n')
    },
    showEvidence(item) {
      this.evidenceItem = item
      this.evidenceDialogVisible = true
    },
    openCertImage() {
      if (this.certInfo && this.certInfo.imgSrc) {
        let src = this.certInfo.imgSrc
        if (src.startsWith('ipfs://')) src = 'http://127.0.0.1:8080/ipfs/' + src.replace('ipfs://', '')
        this.imgPreviewSrc = src
        this.imgDialogVisible = true
      } else {
        this.$message.warning('无证照图片')
      }
    },
    verifyImage() {
      this.$message.info('计算图片SHA-256并与链上指纹比对（需链码存储指纹）')
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
.timeline-only-container { padding: 25px 40px; max-width: 1100px; margin: 0 auto; }
.search-bar { text-align: center; margin: 20px 0; }
.cert-summary { margin-bottom: 20px; }
.summary-row { display: flex; align-items: center; gap: 16px; margin-bottom: 8px; flex-wrap: wrap; }
.history-card { margin-bottom: 8px; }
.history-header { display: flex; justify-content: space-between; margin-bottom: 12px; }
.history-details { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.detail-item { display: flex; flex-direction: column; }
.detail-item.full-width { grid-column: 1 / -1; }
.detail-label { font-size: 12px; color: #909399; margin-bottom: 4px; }
.detail-value { font-size: 13px; color: #303133; word-break: break-all; }
.tx-hash { font-family: monospace; font-size: 12px; background: #f5f7fa; padding: 2px 6px; border-radius: 4px; display: inline-flex; align-items: center; gap: 8px; }
.status-badge { display: flex; align-items: center; gap: 8px; }
.evid-btn { text-align: right; margin-top: 8px; }
.empty-card { text-align: center; padding: 40px; }
.empty-placeholder i { font-size: 48px; color: #c0c4cc; }
</style>
