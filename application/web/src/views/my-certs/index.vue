<template>
  <div class="my-certs-container">
    <h3>我的证照</h3>
    <p style="text-align:center;color:#666">共 {{ certList.length }} 条证照记录</p>

    <el-table v-loading="loading" :data="certList" border stripe style="margin-top:15px">
      <el-table-column prop="traceability_code" label="溯源码" width="180" />
      <el-table-column label="证照类型" width="110">
        <template slot-scope="scope">{{ (scope.row.farmer_input || {}).fa_fruitName }}</template>
      </el-table-column>
      <el-table-column label="证照编号" width="200">
        <template slot-scope="scope">{{ (scope.row.farmer_input || {}).fa_origin }}</template>
      </el-table-column>
      <el-table-column label="上链时间" width="170">
        <template slot-scope="scope">{{ (scope.row.farmer_input || {}).fa_timestamp }}</template>
      </el-table-column>

      <!-- 审核状态 -->
      <el-table-column label="审核状态" width="110">
        <template slot-scope="scope">
          <el-tag v-if="(scope.row.factory_input || {}).fac_contactNumber === '通过'" type="success" size="small">审核通过</el-tag>
          <el-tag v-else-if="(scope.row.factory_input || {}).fac_contactNumber === '驳回'" type="danger" size="small">审核驳回</el-tag>
          <el-tag v-else-if="(scope.row.factory_input || {}).fac_txid" type="warning" size="small">审核中</el-tag>
          <el-tag v-else type="info" size="small">待审核</el-tag>
        </template>
      </el-table-column>

      <!-- 备案状态 -->
      <el-table-column label="备案状态" width="90">
        <template slot-scope="scope">
          <el-tag v-if="(scope.row.driver_input || {}).dr_txid" type="success" size="small">已备案</el-tag>
          <el-tag v-else type="info" size="small">未备案</el-tag>
        </template>
      </el-table-column>

      <!-- 核验状态 -->
      <el-table-column label="核验状态" width="90">
        <template slot-scope="scope">
          <el-tag v-if="(scope.row.shop_input || {}).sh_txid" type="success" size="small">已核验</el-tag>
          <el-tag v-else type="info" size="small">未核验</el-tag>
        </template>
      </el-table-column>

      <!-- 提醒 -->
      <el-table-column label="提醒" width="140">
        <template slot-scope="scope">
          <el-tag v-if="(scope.row.factory_input || {}).fac_contactNumber === '驳回'" type="danger" size="small" effect="dark">
            <i class="el-icon-warning" /> 审核被驳回
          </el-tag>
          <span v-else-if="(scope.row.factory_input || {}).fac_contactNumber === '通过' && !(scope.row.driver_input || {}).dr_txid" style="color:#e6a23c">
            <i class="el-icon-info" /> 等待企业使用
          </span>
          <span v-else-if="(scope.row.shop_input || {}).sh_txid" style="color:#67c23a">
            <i class="el-icon-circle-check" /> 全流程完成
          </span>
          <span v-else style="color:#999">--</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getFruitList } from '@/api/trace'

export default {
  data() {
    return { certList: [], loading: false }
  },
  mounted() { this.loadCerts() },
  methods: {
    loadCerts() {
      this.loading = true
      getFruitList().then(res => {
        if (res.code === 200 && res.data) {
          this.certList = JSON.parse(res.data)
        }
      }).finally(() => { this.loading = false })
    }
  }
}
</script>

<style scoped>
.my-certs-container { padding: 25px 40px; max-width: 1200px; margin: 0 auto; }
h3 { text-align: center; margin-bottom: 8px; }
</style>
