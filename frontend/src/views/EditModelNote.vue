<template>
  <div class="edit-model-note">
    <h2>编辑模型备注</h2>
    <div class="edit-form">
      <el-form :model="form" label-width="80px">
        <el-form-item label="模型名称">
          <span>{{ modelName }}</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="form.note"
            type="textarea"
            :rows="4"
            placeholder="请输入模型备注"
          />
        </el-form-item>
        <el-form-item label="群组ID">
          <el-input
            v-model="form.openChatId"
            placeholder="请输入飞书群组ID"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveNote">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'EditModelNote',
  data() {
    return {
      form: {
        note: '',
        openChatId: ''
      }
    }
  },
  computed: {
    modelName() {
      return this.$route.params.modelName
    }
  },
  methods: {
    async saveNote() {
      try {
        await axios.post('/models/note', {
          model_name: this.modelName,
          note: this.form.note,
          open_chat_id: this.form.openChatId
        })
        this.$message.success('保存成功')
        this.goBack()
      } catch (error) {
        this.$message.error('保存失败：' + error.message)
      }
    },
    goBack() {
      this.$router.push('/models')
    }
  }
}
</script>

<style scoped>
.edit-model-note {
  padding: 20px;
}
.edit-form {
  max-width: 600px;
  margin: 20px auto;
}
</style> 