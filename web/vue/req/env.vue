<template>
  <div>
    <Table border :columns="envTableColumns" :data="envTableData" :no-data-text="$t('goman.general.noData')" :no-filtered-data-text="$t('goman.general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import API from "../../ts/api";
// import _ from 'lodash'
// import moment from 'moment'

@Component
export default class AssertList extends Vue {
  @Prop() isEnvs: boolean
  @Prop() isSending: boolean = false
  @Prop() envListData: Req.EnvListModel[]

  // depDataTableData: Req.DepDataListModel[] = []

  depDataTableData: string[] = []

  envTableData: any[] = []
  envTableColumns: any[] = [
    {
      title: '关联产品',
      key: 'product',
      width: 1120,
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
              props: {
                value: params.row.product,
                transfer: true,
                filterable: true,
                data: this.getOptions('product'),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'product', value)
                }
              }
            },
            this.depDataTableData.map((item) => {
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

  productOptions: string[] = []
  hostIpOptions: string[] = []
  protocolOptions: string[] = ["http", "https"]

  created() {
    this.onLocale();
    this.getEnvList();
  }

  async getEnvData(val) {
    let result = await API.get<Req.ResponseModel[]>('/envList/'+val)
    console.log("getEnvData")
    console.log(result)
    if (result.data) {
      this.depDataTableData = result.data.map(item => (item['product']));
    }
  }

  async getEnvList() {
    let result = await API.get<Req.ResponseModel[]>('/envList')
    if (result.data) {
      this.depDataTableData = result.data.map(item => (item['product']));
    }

    this.envTableColumns =  [
      {
        title: '关联产品',
        key: 'product',
        width: 1120,
        render: (h: CreateElement, params: any) => {
          return h('AutoComplete', {
                props: {
                  value: params.row.product,
                  transfer: true,
                  filterable: true,
                  data: this.getOptions('product'),
                  'filter-method': this.onFilterMethod
                },
                on: {
                  'on-change': (value: string) => {
                    this.onChange(params.index, 'product', value)
                  }
                }
              },
              this.depDataTableData.map((item) => {
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
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.envTableColumns[0].title = this.$t('envList.product')
  }

  mounted() {
    this.syncTableData()
    // this.onRunListDataChange()
  }

  @Watch('envListData')
  onListDataChange() {
    this.syncTableData()
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.envListData.push({
      product: '',
      // hostIp: '',
      // protocol: 'http'
    })
    this.$emit('onChange')
    this.syncTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.envListData.splice(index, 1)
    this.$emit('onChange')

    this.syncTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'product') {
      this.envListData[index].product = value
    } //else if (name === 'hostIp') {
    //   this.envListData[index].hostIp = value
    // }  else if (name === 'protocol') {
    //   this.envListData[index].protocol = value
    // }
    this.$emit('onChange')
  }

  getOptions(type: 'product' | 'hostIp' | 'protocol') {
    if (type==='product') {
      return this.productOptions
    } else if (type==='hostIp') {
      return this.hostIpOptions
    } else if (type==='protocol') {
      return this.protocolOptions
    }
  }

  syncTableData() {
    //全量复制
    this.envTableData = this.envListData.map((m): any => {
      return {
        product: m.product,
        // hostIp: m.hostIp,
        // protocol: m.protocol
      }
    })
  }

}
</script>

