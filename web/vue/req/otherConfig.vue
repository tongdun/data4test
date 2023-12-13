<template>
  <div>
    <Table border :columns="otherConfigTableColumns" :data="otherConfigTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import _ from 'lodash'
// import moment from 'moment'

@Component
export default class OtherConfigList extends Vue {
  @Prop() isSending: boolean = false
  @Prop() contentType: string
  @Prop() otherConfigListData: Req.OtherConfigListModel[]

  yOrNLable: string[] = ["yes", "no"]
  otherConfigTableData: any[] = []
  otherConfigTableColumns: any[] = [
    {
      title: '',
      key: 'version',
      width: 100,
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
              props: {
                value: params.row.version,
                transfer: true,
                disabled: true,
                data: this.getOptions('version'),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'version', value)
                }
              }
            },)
      }
    },
    {
      title: '',
      width: 300,
      key: 'apiId',
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
          props: {
            value: params.row.apiId,
            transfer: true,
            disabled: true,
            data: this.getOptions('apiId'),
            'filter-method': this.onFilterMethod
          },
          on: {
            'on-change': (value: string) => {
              this.onChange(params.index, 'apiId', value)
            }
          }
        })
      }
    },
    {
      title: '',
      key: 'isParallel',
      width: 120,
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.isParallel,
                transfer: true,
                data: this.getOptions("isParallel"),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.isParallel, 'value', value)
                }
              }
            },
            this.yOrNLable.map((item) => {
              return h('Option', {
                props: {value: item}
              }, item)
            })
        )
      }
    },
    {
      title: '',
      key: 'isUseEnvConfig',
      width: 240,
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.isUseEnvConfig,
                transfer: true,
                data: this.getOptions("isUseEnvConfig"),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.isUseEnvConfig, 'value', value)
                }
              }
            },
            this.yOrNLable.map((item) => {
              return h('Option', {
                props: {value: item}
              }, item)
            })
        )
      }
    },
    {
      width: 80,
      renderHeader: (h: CreateElement, params: any) => {
        return h('Button', {
          props: {
            disabled: this.isSending,
            size: 'small',
            type: 'primary',
            icon: 'plus-round'
          },
          on: {
            click: () => {
              this.onAdd()
            }
          }
        })
      },
      render: (h: CreateElement, params: any) => {
        return h('Button', {
          props: {
            disabled: this.isSending,
            size: 'small',
            icon: 'close-round'
          },
          on: {
            click: () => {
              this.onDel(params.index)
            }
          }
        })
      }
    }
  ]

  versionOptions: string[] = []
  apiIdOptions: string[] = []
  isParallelOptions: string[] = []
  isUseEnvConfigOptions: string[] = []

  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.otherConfigTableColumns[0].title = this.$t('otherConfigList.version')
    this.otherConfigTableColumns[1].title = this.$t('otherConfigList.apiId')
    this.otherConfigTableColumns[2].title = this.$t('otherConfigList.isParallel')
    this.otherConfigTableColumns[3].title = this.$t('otherConfigList.isUseEnvConfig')

  }

  mounted() {
    this.syncTableData()
  }

  @Watch('otherConfigListData')
  onListDataChange() {
    this.syncTableData()
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.otherConfigListData.push({
      version: '',
      apiId: '',
      isParallel: 'yes',
      isUseEnvConfig: 'yes'
    })
    this.$emit('onChange')
    this.syncTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.otherConfigListData.splice(index, 1)
    this.$emit('onChange')

    this.syncTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'version') {
      this.otherConfigListData[index].version = value
    } else if (name === 'apiId') {
      this.otherConfigListData[index].apiId = value
    }  else if (name === 'isParallel') {
      this.otherConfigListData[index].isParallel = value
    } else if (name === 'isUseEnvConfig') {
      this.otherConfigListData[index].isUseEnvConfig = value
    }
    this.$emit('onChange')
  }

  getOptions(type: 'version' | 'apiId' | 'isParallel' | 'isUseEnvConfig') {
    if (type==='version') {
      return this.versionOptions
    } else if (type==='apiId') {
      return this.apiIdOptions
    } else if (type==='isParallel') {
      return this.isParallelOptions
    } else if (type==='isUseEnvConfig') {
      return this.isUseEnvConfigOptions
    }

  }

  syncTableData() {
    //全量复制
    this.otherConfigTableData = this.otherConfigListData.map((m): any => {
      return {
        version: m.version,
        apiId: m.apiId,
        isParallel: m.isParallel,
        isUseEnvConfig: m.isUseEnvConfig
      }
    })
  }
}
</script>

