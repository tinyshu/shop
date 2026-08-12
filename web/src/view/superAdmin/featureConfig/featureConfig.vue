<template>
  <div class="feature-config">
    <div class="gva-search-box">
      <el-alert
        title="行业/功能开关（sys_config · group=feature）。修改后下次业务请求生效；支付等 yaml 配置不在此页。"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <div class="gva-table-box" v-loading="loading">
      <el-alert
        v-if="missingKeys.length"
        class="mb-4"
        :title="`缺少配置键：${missingKeys.join('、')}。请先执行 sql/migrations/20260811_industry_config_m0_1.sql`"
        type="warning"
        :closable="false"
        show-icon
      />

      <el-form label-width="160px" style="max-width: 640px">
        <el-form-item label="用户审核">
          <div class="row">
            <el-switch
              v-model="form.userAudit"
              :disabled="!rows.userAudit || saving"
              @change="(v) => onBoolChange('userAudit', v)"
            />
            <span class="hint">{{ form.userAudit ? '开启：未审核用户不可下单（小 B）' : '关闭：B2C 免审（默认）' }}</span>
          </div>
          <div class="key">feature.user_audit</div>
        </el-form-item>

        <el-form-item label="月结能力">
          <div class="row">
            <el-switch
              v-model="form.settleMonth"
              :disabled="!rows.settleMonth || saving"
              @change="(v) => onBoolChange('settleMonth', v)"
            />
            <span class="hint">占位开关；履约侧业务后续切片才真正拦截/展示</span>
          </div>
          <div class="key">feature.settle_month</div>
        </el-form-item>

        <el-form-item label="物流模式">
          <div class="row">
            <el-radio-group
              v-model="form.courierMode"
              :disabled="!rows.courierMode || saving"
              @change="onCourierChange"
            >
              <el-radio label="courier">快递单号</el-radio>
              <el-radio label="delivery">城配配送员</el-radio>
            </el-radio-group>
          </div>
          <div class="hint">占位开关；发货字段联动见 M4</div>
          <div class="key">feature.courier_mode</div>
        </el-form-item>
      </el-form>

      <el-button icon="refresh" :loading="loading" @click="load">刷新</el-button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FeatureConfig'
}
</script>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getSysConfigList, updateSysConfig } from '@/api/sysConfig'

const KEYS = {
  userAudit: 'feature.user_audit',
  settleMonth: 'feature.settle_month',
  courierMode: 'feature.courier_mode'
}

const loading = ref(false)
const saving = ref(false)
const rows = reactive({
  userAudit: null,
  settleMonth: null,
  courierMode: null
})
const form = reactive({
  userAudit: false,
  settleMonth: false,
  courierMode: 'courier'
})

const missingKeys = computed(() => {
  const miss = []
  if (!rows.userAudit) miss.push(KEYS.userAudit)
  if (!rows.settleMonth) miss.push(KEYS.settleMonth)
  if (!rows.courierMode) miss.push(KEYS.courierMode)
  return miss
})

const parseBool = (value) => {
  const v = String(value || '').trim().toLowerCase()
  return v === '1' || v === 'true' || v === 'on'
}

const findByName = (list, name) => (list || []).find((item) => item.name === name) || null

const load = async() => {
  loading.value = true
  try {
    const res = await getSysConfigList({ page: 1, pageSize: 50, groupType: 'feature' })
    if (res.code !== 0) {
      ElMessage.error(res.msg || '加载失败')
      return
    }
    const list = res.data?.list || []
    rows.userAudit = findByName(list, KEYS.userAudit)
    rows.settleMonth = findByName(list, KEYS.settleMonth)
    rows.courierMode = findByName(list, KEYS.courierMode)
    form.userAudit = rows.userAudit ? parseBool(rows.userAudit.value) : false
    form.settleMonth = rows.settleMonth ? parseBool(rows.settleMonth.value) : false
    form.courierMode = rows.courierMode?.value === 'delivery' ? 'delivery' : 'courier'
  } finally {
    loading.value = false
  }
}

const saveRow = async(row, nextValue, rollback) => {
  if (!row) return
  saving.value = true
  try {
    const payload = {
      ...row,
      value: nextValue,
      status: 1
    }
    const res = await updateSysConfig(payload)
    if (res.code !== 0) {
      rollback()
      ElMessage.error(res.msg || '保存失败')
      return
    }
    row.value = nextValue
    row.status = 1
    ElMessage.success('已保存')
  } catch (e) {
    rollback()
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const onBoolChange = (field, enabled) => {
  const row = rows[field]
  saveRow(row, enabled ? '1' : '0', () => {
    form[field] = !enabled
  })
}

const onCourierChange = (mode) => {
  const row = rows.courierMode
  const before = row?.value === 'delivery' ? 'delivery' : 'courier'
  saveRow(row, mode, () => {
    form.courierMode = before
  })
}

onMounted(load)
</script>

<style scoped>
.mb-4 {
  margin-bottom: 16px;
}
.row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}
.hint {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.key {
  margin-top: 4px;
  color: var(--el-text-color-placeholder);
  font-size: 12px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
}
</style>
