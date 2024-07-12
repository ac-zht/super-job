<template>
  <el-container>
    <task-sidebar></task-sidebar>
    <el-main>
      <el-breadcrumb separator-class="el-icon-arrow-right" style="margin-bottom:20px">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/task' }">任务管理</el-breadcrumb-item>
        <el-breadcrumb-item>编辑</el-breadcrumb-item>
      </el-breadcrumb>
      <el-form ref="form" class="page-form" :model="form" :rules="formRules" label-width="180px">
        <el-input v-model="form.id" type="hidden"></el-input>
        <el-row>
          <el-col :span="15">
            <el-form-item label="任务名称" prop="name">
              <el-input v-model.trim="form.name"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="15" style="margin-left: 180px;">
            <i class="el-icon-time" style="color: #909399;"></i>
            <el-button
              style="color: #909399;"
              class="box-shadow-not"
              type="text"
              v-for="(item, index) in specOptions"
              :key="index"
              @click="specSelect(item.value)"
            >{{ item.label }}
            </el-button
            >
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="15">
            <el-form-item label="crontab表达式" prop="expression">
              <el-input
                v-model.trim="form.expression"
                placeholder="秒 分 时 天 月 周"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8">
            <el-form-item label="执行方式">
              <el-select v-model.trim="form.protocol">
                <el-option
                  v-for="item in protocolList"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8" v-if="form.protocol === 1">
            <el-form-item label="请求方法">
              <el-select key="http-method" v-model.trim="form.http_method">
                <el-option
                  v-for="item in httpMethods"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8" v-else>
            <el-form-item label="执行器">
              <el-select v-model.trim="form.exec_id">
                <el-option label="请选择" value=""></el-option>
                <el-option
                  v-for="item in executors"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8">
            <el-form-item label="单实例运行">
              <span slot="label">
                单实例运行
                <el-tooltip placement="top">
                  <div slot="content">
                    单实例运行,
                    前次任务未执行完成，下次任务调度时间到了是否要执行,
                    即是否允许多进程执行同一任务
                  </div>
                  <i class="el-icon-question"></i>
                </el-tooltip>
              </span>
              <el-select v-model.trim="form.multi">
                <el-option
                  v-for="item in runStatusList"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8" v-if="form.protocol === 2">
            <el-form-item label="JobHandler" prop="executor_handler">
              <el-input v-model.trim="form.executor_handler"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="15" v-else>
            <el-form-item label="命令" prop="command">
              <el-input
                type="textarea"
                :rows="5"
                size="medium"
                width="100"
                :placeholder="commandPlaceholder"
                v-model="form.command">
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="15">
            <el-form-item label="任务超时时间" prop="timeout">
              <span slot="label">
                任务超时时间
                <el-tooltip placement="top">
                  <div slot="content">
                    任务执行超时强制结束, 取值0-86400(秒), 默认0, 不限制
                  </div>
                  <i class="el-icon-question"></i>
                </el-tooltip>
              </span>
              <el-input v-model.number.trim="form.timeout"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="15">
            <el-form-item label="任务失败重试次数" prop="retry_times">
              <el-input
                v-model.number.trim="form.retry_times"
                placeholder="0 - 10, 默认0，不重试"
              ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="15">
            <el-form-item label="任务失败重试间隔时间" prop="retry_interval">
              <el-input
                v-model.number.trim="form.retry_interval"
                placeholder="0 - 3600 (秒), 默认0，执行系统默认策略"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="15">
            <el-form-item label="状态">
            <span slot="label">
                状态
                <el-tooltip placement="top">
                  <div slot="content">
                    开关禁用表示任务正在执行
                  </div>
                  <i class="el-icon-question"></i>
                </el-tooltip>
              </span>
              <el-switch
                v-model="logic.status"
                :active-value="1"
                :inactive-vlaue="0"
                :disabled="form.status===2"
              >
              </el-switch>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8">
            <el-form-item label="任务通知">
              <el-select v-model.trim="form.notify_status">
                <el-option
                  v-for="item in notifyStatusList"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8" v-if="form.notify_status !== 0">
            <el-form-item label="通知类型">
              <el-select v-model.trim="form.notify_type">
                <el-option
                  v-for="item in notifyTypes"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col
            :span="8"
            v-if="form.notify_status !== 0 && form.notify_type === 1"
          >
            <el-form-item label="接收用户">
              <el-select
                key="notify-mail"
                v-model="selectedMailNotifyIds"
                filterable
                multiple
                placeholder="请选择"
              >
                <el-option
                  v-for="item in mailUsers"
                  :key="item.id"
                  :label="item.username"
                  :value="item.id"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>

          <el-col
            :span="8"
            v-if="form.notify_status !== 0 && form.notify_type === 2"
          >
            <el-form-item label="发送Channel">
              <el-select
                key="notify-slack"
                v-model="selectedSlackNotifyIds"
                filterable
                multiple
                placeholder="请选择"
              >
                <el-option
                  v-for="item in slackChannels"
                  :key="item.id"
                  :label="item.name"
                  selected="true"
                  :value="item.id"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="form.notify_status === 3">
          <el-col :span="15">
            <el-form-item label="任务执行输出关键字" prop="notify_keyword">
              <el-input
                v-model.trim="form.notify_keyword"
                placeholder="任务执行输出中包含此关键字将触发通知"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="15">
            <el-form-item label="备注">
              <el-input
                type="textarea"
                :rows="3"
                size="medium"
                width="100"
                v-model="form.cfg"
              >
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" @click="submit">保存</el-button>
          <el-button @click="cancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
import taskSidebar from './sidebar'
import taskService from '../../api/task'
import notificationService from '../../api/notification'

export default {
  name: 'task-edit',
  data() {
    return {
      logic: {
        status: 0
      },
      form: {
        id: '',
        name: '',
        expression: '',
        protocol: 2,
        http_method: 1,
        executor_handler: '',
        command: '',
        exec_id: '',
        timeout: 0,
        multi: 2,
        notify_status: 0,
        notify_type: 1,
        notify_receiver_id: '',
        notify_keyword: '',
        retry_times: 0,
        retry_interval: 0,
        status: 1,
        cfg: ''
      },
      formRules: {
        name: [
          {required: true, message: '请输入任务名称', trigger: 'blur'}
        ],
        expression: [
          {required: true, message: '请输入crontab表达式', trigger: 'blur'}
        ],
        executor_handler: [
          {required: true, message: '请输入方法', trigger: 'blur'}
        ],
        command: [
          {required: true, message: '请输入命令', trigger: 'blur'}
        ],
        timeout: [
          {
            type: 'number',
            required: true,
            message: '请输入有效的任务超时时间',
            trigger: 'blur'
          }
        ],
        retry_times: [
          {
            type: 'number',
            required: true,
            message: '请输入有效的任务执行失败重试次数',
            trigger: 'blur'
          }
        ],
        retry_interval: [
          {
            type: 'number',
            required: true,
            message: '请输入有效的任务执行失败，重试间隔时间',
            trigger: 'blur'
          }
        ],
        notify_keyword: [
          {
            required: true,
            message: '请输入要匹配的任务执行输出关键字',
            trigger: 'blur'
          }
        ]
      },
      httpMethods: [
        {
          value: 1,
          label: 'get'
        },
        {
          value: 2,
          label: 'post'
        }
      ],
      protocolList: [
        {
          value: 1,
          label: 'http'
        },
        {
          value: 2,
          label: 'rpc'
        },
        {
          value: 3,
          label: 'shell'
        }
      ],
      runStatusList: [
        {
          value: 1,
          label: '是'
        },
        {
          value: 0,
          label: '否'
        }
      ],
      notifyStatusList: [
        {
          value: 0,
          label: '不通知'
        },
        {
          value: 1,
          label: '失败通知'
        },
        {
          value: 2,
          label: '总是通知'
        },
        {
          value: 3,
          label: '关键字匹配通知'
        }
      ],
      notifyTypes: [
        {
          value: 1,
          label: '邮件'
        },
        {
          value: 2,
          label: 'Slack'
        },
        {
          value: 3,
          label: 'WebHook'
        }
      ],
      executors: [],
      mailUsers: [],
      slackChannels: [],
      selectedMailNotifyIds: [],
      selectedSlackNotifyIds: [],
      specOptions: [
        {
          value: '0 * * * * *',
          label: '每分钟'
        },
        {
          value: '0 */5 * * * *',
          label: '每5分钟'
        },
        {
          value: '0 0 * * * *',
          label: '每1小时整'
        },
        {
          value: '0 0 0 * * *',
          label: '每天0点整'
        },
        {
          value: '0 0 0 * * 1',
          label: '每周周一'
        },
        {
          value: '0 0 0 1 * *',
          label: '每月第1天'
        }
      ]
    }
  },
  computed: {
    commandPlaceholder() {
      if (this.form.protocol === 1) {
        return '请输入URL地址'
      }
      if (this.form.protocol === 3) {
        return '请输入shell命令'
      }
      // return '请输入shell命令'
    }
  },
  components: {taskSidebar},
  created() {
    const id = this.$route.params.id

    taskService.detail(id, (taskData, executors) => {
      if (id && !taskData) {
        this.$message.error('数据不存在')
        this.cancel()
        return
      }
      this.executors = executors || []
      if (!taskData) {
        return
      }
      this.form.id = taskData.id
      this.form.name = taskData.name
      this.form.expression = taskData.expression
      this.form.protocol = taskData.protocol
      if (this.form.protocol === 2) {
        this.form.exec_id = taskData.exec_id
        this.form.executor_handler = taskData.executor_handler
      } else {
        this.form.command = taskData.command
      }
      if (taskData.http_method) {
        this.form.http_method = taskData.http_method
      }
      this.form.timeout = taskData.timeout
      this.form.multi = taskData.multi ? 1 : 0
      this.form.notify_keyword = taskData.notify_keyword
      this.form.notify_status = taskData.notify_status
      this.form.status = taskData.status
      this.logic.status = (taskData.status > 0) ? 1 : 0
      this.form.notify_receiver_id = taskData.notify_receiver_id
      this.form.notify_type = taskData.notify_type
      this.form.retry_times = taskData.retry_times
      this.form.retry_interval = taskData.retry_interval
      this.form.cfg = taskData.cfg

      if (this.form.notify_status > 0) {
        const notifyReceiverIds = this.form.notify_receiver_id.split(',')
        if (this.form.notify_type === 1) {
          notifyReceiverIds.forEach((v) => {
            this.selectedMailNotifyIds.push(parseInt(v))
          })
        } else if (this.form.notify_type === 2) {
          notifyReceiverIds.forEach((v) => {
            this.selectedSlackNotifyIds.push(parseInt(v))
          })
        }
      }
    })

    notificationService.mail((data) => {
      this.mailUsers = data.mail_users
    })

    notificationService.slack((data) => {
      this.slackChannels = data.channels
    })
  },
  methods: {
    submit() {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        if (this.form.notify_status > 0) {
          if (
            this.form.notify_type === 1 &&
            this.selectedMailNotifyIds.length === 0
          ) {
            this.$message.error('请选择邮件接收用户')
            return false
          }
          if (
            this.form.notify_type === 2 &&
            this.selectedSlackNotifyIds.length === 0
          ) {
            this.$message.error('请选择Slack Channel')
            return false
          }
        }

        this.save()
      })
    },
    save() {
      if (this.form.notify_status > 0 && this.form.notify_type === 1) {
        this.form.notify_receiver_id = this.selectedMailNotifyIds.join(',')
      }
      if (this.form.notify_status > 0 && this.form.notify_type === 2) {
        this.form.notify_receiver_id = this.selectedSlackNotifyIds.join(',')
      }
      if (this.form.status !== 2) {
        this.form.status = this.logic.status
      }
      taskService.update(this.form, () => {
        this.$router.push('/task')
      })
    },
    cancel() {
      this.$router.push('/task')
    },
    specSelect(expression) {
      this.form.expression = expression
    }
  }
}
</script>
