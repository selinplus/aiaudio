<template>
    <v-app>
    <v-content>
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex justify-center>
            <v-card class="elevation-12">
              <v-card-title  primary-title v-if="status">
                 <v-icon color="red lighten-1" xLarge flat>record_voice_over</v-icon>
              </v-card-title>
              <v-card-text align-center>                
                <v-btn flat large outline @click="start" fluid>
                   <v-icon>speaker_phone</v-icon>                  
                 </v-btn>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn flat small disabled>测试录音</v-btn>
              </v-card-actions>
            </v-card>
          </v-flex>          
          <v-snackbar
            v-model="snackbar"
            :color="color"
            :multi-line="mode === 'multi-line'"
            :timeout="timeout"
            :vertical="mode === 'vertical'"
          >
            {{ text }}
            <v-btn dark flat @click="snackbar = false">Close</v-btn>
          </v-snackbar>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>
<script>
export default {
  data:() =>({
    snackbar: false,
    color: "orange",
    mode: "",
    timeout: 3000,
    text: "",
    tip:"录制",
    status:false
  }),
  mounted(){
    window.wails.Events.On('NEWS',(n,t) => {
      this.snackbar = true
      this.text = "[" + t + "]" + n      
    })
  },
  methods:{
    start(){
      if(this.status){
        window.wails.Events.Emit('end')
      }else{
        window.wails.Events.Emit('start')
      }      
      this.status =!this.status
    },
    end(){
      window.wails.Events.Emit('end')
      this.status =!this.status
    },    
  }
}
</script>