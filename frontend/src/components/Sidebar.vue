        <el-menu-item 
          index="/resources" 
          v-if="['admin', 'operator'].includes(userStore.user.role)"
        >
          <i class="el-icon-monitor"></i>
          <span>资源管理</span>
        </el-menu-item> 