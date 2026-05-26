<template>
  <div class="tech-container">
    <h3>技术支撑实体 — 证照核验</h3>
    <p style="text-align:center;color:#666">当前角色：技术支撑实体 | 待核验：{{ pendingList.length }} 条</p>

    <!-- 统计卡片区 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="12">
        <el-card
          class="stat-card"
          :class="{ active: activeTab === 'pending' }"
          shadow="hover"
          @click.native="activeTab = 'pending'"
        >
          <div class="stat-icon"><i class="el-icon-warning-outline" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ pendingList.length }}</div>
            <div class="stat-title">待核验证照</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card
          class="stat-card"
          :class="{ active: activeTab === 'verified' }"
          shadow="hover"
          @click.native="activeTab = 'verified'"
        >
          <div class="stat-icon"><i class="el-icon-circle-check" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ verifiedList.length }}</div>
            <div class="stat-title">已核验证照</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待核验证照列表 -->
    <div v-show="activeTab === 'pending'">
      <el-table v-loading="loading" :data="pendingList" border stripe fit>
        <el-table-column prop="traceability_code" label="溯源码" min-width="180" />
        <el-table-column label="证照类型" min-width="110">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_fruitName }}</template>
        </el-table-column>
        <el-table-column label="持证人" min-width="100">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_farmerName }}</template>
        </el-table-column>
        <el-table-column label="审核状态" min-width="90">
          <template slot-scope="scope">
            <el-tag v-if="(scope.row.factory_input || {}).fac_contactNumber === '通过'" type="success">已通过</el-tag>
            <el-tag v-else-if="(scope.row.factory_input || {}).fac_txid" type="danger">已驳回</el-tag>
            <el-tag v-else type="warning">待审核</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="备案状态" min-width="90">
          <template slot-scope="scope">
            <el-tag v-if="(scope.row.driver_input || {}).dr_txid" type="success">已备案</el-tag>
            <el-tag v-else type="info">未备案</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="100">
          <template slot-scope="scope">
            <el-button v-if="canVerify(scope.row)" type="primary" size="small" @click="openVerify(scope.row)">核验</el-button>
            <el-tag v-else-if="(scope.row.factory_input || {}).fac_contactNumber !== '通过'" type="warning" size="small">待审核</el-tag>
            <el-tag v-else type="info" size="small">未备案</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 已核验证照列表 -->
    <div v-show="activeTab === 'verified'">
      <el-card style="margin-top:25px">
        <div slot="header">已核验证照（{{ verifiedList.length }} 条）</div>
        <el-table :data="verifiedList" border stripe fit>
          <el-table-column prop="traceability_code" label="溯源码" min-width="180" />
          <el-table-column label="证照类型" min-width="110">
            <template slot-scope="scope">{{ scope.row.farmer_input.fa_fruitName }}</template>
          </el-table-column>
          <el-table-column label="持证人" min-width="100">
            <template slot-scope="scope">{{ scope.row.farmer_input.fa_farmerName }}</template>
          </el-table-column>
          <el-table-column label="核验实体" min-width="140">
            <template slot-scope="scope">{{ scope.row.shop_input.sh_storeTime }}</template>
          </el-table-column>
          <el-table-column label="核验结果" min-width="90">
            <template slot-scope="scope">
              <el-tag v-if="scope.row.shop_input.sh_shopAddress === '有效'" type="success">有效</el-tag>
              <el-tag v-else type="danger">无效</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="核验时间" min-width="160">
            <template slot-scope="scope">{{ scope.row.shop_input.sh_timestamp }}</template>
          </el-table-column>
          <el-table-column label="操作" min-width="100">
            <template slot-scope="scope">
              <el-button size="small" @click="viewDetail(scope.row)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 核验对话框 -->
    <el-dialog title="证照核验" :visible.sync="dialogVisible" width="500px">
      <el-descriptions v-if="currentCert" :column="1" border>
        <el-descriptions-item label="溯源码">{{ currentCert.traceability_code }}</el-descriptions-item>
        <el-descriptions-item label="证照类型">{{ (currentCert.farmer_input || {}).fa_fruitName }}</el-descriptions-item>
        <el-descriptions-item label="持证人">{{ (currentCert.farmer_input || {}).fa_farmerName }}</el-descriptions-item>
        <el-descriptions-item label="使用企业">{{ (currentCert.driver_input || {}).dr_name }}</el-descriptions-item>
        <el-descriptions-item label="使用目的">{{ (currentCert.driver_input || {}).dr_phone }}</el-descriptions-item>
      </el-descriptions>

      <el-form ref="verifyForm" :model="vfy" :rules="rules" label-width="110px" style="margin-top:20px">
        <el-form-item label="实体名称" prop="entityName">
          <el-select v-model="vfy.entityName" placeholder="请选择" style="width:100%">
            <el-option label="中认CA认证中心" value="中认CA认证中心" />
            <el-option label="国信电子认证中心" value="国信电子认证中心" />
            <el-option label="数安密码服务中心" value="数安密码服务中心" />
            <el-option label="可信时间戳中心" value="可信时间戳中心" />
            <el-option label="华测检测认证" value="华测检测认证" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务类型" prop="serviceType">
          <el-select v-model="vfy.serviceType" placeholder="请选择" style="width:100%">
            <el-option label="CA数字证书认证" value="CA数字证书认证" />
            <el-option label="时间戳服务" value="时间戳服务" />
            <el-option label="电子签名验证" value="电子签名验证" />
            <el-option label="哈希校验服务" value="哈希校验服务" />
            <el-option label="数据加密传输" value="数据加密传输" />
          </el-select>
        </el-form-item>
        <el-form-item label="安全认证等级" prop="securityLevel">
          <el-select v-model="vfy.securityLevel" placeholder="请选择" style="width:100%">
            <el-option label="一级（基础）" value="一级（基础）" />
            <el-option label="二级（增强）" value="二级（增强）" />
            <el-option label="三级（高级）" value="三级（高级）" />
            <el-option label="四级（最高）" value="四级（最高）" />
          </el-select>
        </el-form-item>
      </el-form>

      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="success" @click="doVerify('有效')">核验有效</el-button>
        <el-button type="danger" @click="doVerify('无效')">核验无效</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getAllFruitInfo } from '@/api/trace'
import request from '@/utils/request'

export default {
  data() {
    return {
      allCerts: [],
      loading: false,
      dialogVisible: false,
      currentCert: null,
      vfy: { entityName: '', serviceType: '', securityLevel: '' },
      rules: {
        entityName: [{ required: true, message: '请选择实体' }],
        serviceType: [{ required: true, message: '请选择服务类型' }],
        securityLevel: [{ required: true, message: '请选择安全等级' }]
      },
      activeTab: 'pending'
    }
  },
  computed: {
    pendingList() {
      return this.allCerts.filter(c => {
        const fi2 = c.factory_input || {}
        const di = c.driver_input || {}
        const si = c.shop_input || {}
        return fi2.fac_contactNumber === '通过' && !!di.dr_txid && !si.sh_txid
      })
    },
    verifiedList() { return this.allCerts.filter(c => !!(c.shop_input || {}).sh_txid) }
  },
  mounted() { this.loadCerts() },
  methods: {
    canVerify(row) {
      const fi2 = row.factory_input || {}
      const di = row.driver_input || {}
      return fi2.fac_contactNumber === '通过' && !!di.dr_txid
    },
    loadCerts() {
      this.loading = true
      getAllFruitInfo().then(res => {
        if (res.code === 200 && res.data) this.allCerts = JSON.parse(res.data)
      }).finally(() => { this.loading = false })
    },
    openVerify(cert) {
      this.currentCert = cert
      this.vfy = { entityName: '', serviceType: '', securityLevel: '' }
      this.dialogVisible = true
    },
    doVerify(result) {
      this.$refs.verifyForm.validate(valid => {
        if (!valid) return
        const fd = new FormData()
        fd.append('traceability_code', this.currentCert.traceability_code)
        fd.append('entityName', this.vfy.entityName)
        fd.append('serviceType', this.vfy.serviceType)
        fd.append('securityLevel', this.vfy.securityLevel)
        fd.append('verifyResult', result)
        request({ url: '/techVerify', method: 'post', data: fd }).then(res => {
          if (res.code === 200) {
            this.$message.success('核验' + result + '！')
            this.dialogVisible = false
            this.loadCerts()
          } else {
            this.$message.error(res.message || '核验失败')
          }
        })
      })
    },
    viewDetail(row) {
      this.$router.push({ path: '/trace', query: { code: row.traceability_code }})
    }
  }
}
</script>

<style scoped>
.tech-container { padding: 20px; max-width: 1100px; margin: 0 auto; }
h3 { text-align: center; margin-bottom: 5px; }
.stat-cards {
  margin-bottom: 20px;
}
.stat-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  border-radius: 16px;
}
.stat-card.active {
  border-color: #409eff;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}
.stat-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 10px;
}
.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}
.stat-title {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}
</style>
