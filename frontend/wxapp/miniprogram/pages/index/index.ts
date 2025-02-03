// index.ts
// 获取应用实例
const app = getApp<IAppOption>()
const defaultAvatarUrl = 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'

Page({
  data: {
    motto: 'Hello World',
    userInfo: {
      avatarUrl: defaultAvatarUrl,
      nickName: '',
    },
    hasUserInfo: false,
    canIUseGetUserProfile: wx.canIUse('getUserProfile'),
    canIUseNicknameComp: wx.canIUse('input.type.nickname'),
    isLogin: false,
    loading: false,
    loadingProfile: false,
    profile: null
  },

  onLoad() {
    // 检查登录状态
    const token = wx.getStorageSync('token')
    this.setData({
      isLogin: !!token
    })
  },

  async handleLogin() {
    if (this.data.loading) return
    
    this.setData({ loading: true })
    
    try {
      // 调用微信登录
      const { code } = await wx.login()
      console.log('获取到登录code:', code)
      
      // 发送 code 到后端
      console.log('开始请求后端接口...')
      wx.request({
        url: 'http://192.168.3.36:8080/api/v1/login/wx-mini',
        method: 'POST',
        data: { code },
        header: {
          'content-type': 'application/json'
        },
        success: (res) => {
          console.log('请求成功:', res)
          const { data } = res as any
          
          if (data.code === 200) {
            // 保存 token
            wx.setStorageSync('token', data.data.token)
            // 保存用户信息
            wx.setStorageSync('userInfo', data.data.user)
            
            this.setData({
              isLogin: true
            })

            wx.showToast({
              title: '登录成功',
              icon: 'success'
            })
          } else {
            wx.showToast({
              title: data.message || '登录失败',
              icon: 'error'
            })
          }
        },
        fail: (err) => {
          console.error('请求失败:', err)
          wx.showToast({
            title: '网络请求失败',
            icon: 'error'
          })
        },
        complete: () => {
          this.setData({ loading: false })
        }
      })
    } catch (err: any) {
      console.error('登录失败，详细错误:', err)
      wx.showToast({
        title: err.message || '登录失败',
        icon: 'error'
      })
      this.setData({ loading: false })
    }
  },

  getProfile() {
    if (this.data.loadingProfile) return
    
    this.setData({ loadingProfile: true })
    
    // 获取本地存储的token
    const token = wx.getStorageSync('token')
    if (!token) {
      wx.showToast({
        title: '请先登录',
        icon: 'error'
      })
      this.setData({ 
        loadingProfile: false,
        isLogin: false 
      })
      return
    }

    // 请求用户信息
    wx.request({
      url: 'http://192.168.3.36:8080/api/v1/user/profile',
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
        'content-type': 'application/json'
      },
      success: (res) => {
        console.log('获取用户信息成功:', res)
        const { data } = res as any
        
        if (data.code === 200) {
          this.setData({
            profile: data.data
          })
        } else if (data.code === 401) {
          // token 过期或无效
          wx.removeStorageSync('token')
          wx.removeStorageSync('userInfo')
          this.setData({ isLogin: false })
          wx.showToast({
            title: '登录已过期',
            icon: 'error'
          })
        } else {
          wx.showToast({
            title: data.message || '获取信息失败',
            icon: 'error'
          })
        }
      },
      fail: (err) => {
        console.error('获取用户信息失败:', err)
        wx.showToast({
          title: '网络请求失败',
          icon: 'error'
        })
      },
      complete: () => {
        this.setData({ loadingProfile: false })
      }
    })
  },

  methods: {
    // 事件处理函数
    bindViewTap() {
      wx.navigateTo({
        url: '../logs/logs',
      })
    },
    onChooseAvatar(e: any) {
      const { avatarUrl } = e.detail
      const { nickName } = this.data.userInfo
      this.setData({
        "userInfo.avatarUrl": avatarUrl,
        hasUserInfo: nickName && avatarUrl && avatarUrl !== defaultAvatarUrl,
      })
    },
    onInputChange(e: any) {
      const nickName = e.detail.value
      const { avatarUrl } = this.data.userInfo
      this.setData({
        "userInfo.nickName": nickName,
        hasUserInfo: nickName && avatarUrl && avatarUrl !== defaultAvatarUrl,
      })
    },
    getUserProfile() {
      // 推荐使用wx.getUserProfile获取用户信息，开发者每次通过该接口获取用户个人信息均需用户确认，开发者妥善保管用户快速填写的头像昵称，避免重复弹窗
      wx.getUserProfile({
        desc: '展示用户信息', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
        success: (res) => {
          console.log(res)
          this.setData({
            userInfo: res.userInfo,
            hasUserInfo: true
          })
        }
      })
    },
  },
})
