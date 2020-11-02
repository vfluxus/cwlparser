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
      I=$(inputs.input_sam) O=$(inputs.output_fixmate_bam) M=$(inputs.output_fixmate_bam).txt VALIDATION_STRINGENCY=SILENT OPTICAL_DUPLICATE_PIXEL_DISTANCE=2500 ASSUME_SORT_ORDER=queryname CREATE_MD5_FILE=false

baseCommand: [picard, MarkDuplicates]
