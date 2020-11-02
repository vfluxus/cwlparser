cwlVersion: v1.0

class: CommandLineTool
id: sort_bam
requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement 
    ramMin: 2000

inputs:
  input_bam: 
    type: File
  output_sorted_bam:
    type: string 
  
outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.bam'
    secondaryFiles: [^.bai, ^.bam.md5]

arguments:
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      I=$(inputs.input_bam) O=$(inputs.output_sorted_bam) SO=coordinate CREATE_INDEX=true CREATE_MD5_FILE=false

baseCommand: [picard, SortSam]
