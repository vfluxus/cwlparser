cwlVersion: v1.0

class: CommandLineTool
id: fix_mate
requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement 
    ramMin: 2000

inputs:
  input_bam: 
    type: File
  output_md_bam:
    type: string 
  output_report_txt:
    type: string
  
outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.bam$'
    secondaryFiles: [^.bai, ^.bam.md5]

arguments:
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      I=$(inputs.input_bam) O=$(input.output_md_bam)  M=$(inputs.output_report_txt) ASO=queryname TMP_DIR=tmp

baseCommand: [picard, MarkDuplicates]
