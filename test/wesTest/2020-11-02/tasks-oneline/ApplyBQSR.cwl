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
  indexed_reference_fasta:
    type: File 
    secondaryFiles: [.64.amb, .64.ann, .64.bwt, .64.pac, .64.sa, .64.alt, .dict, .fai]
  input_recalibrate_bqsr:
    type: File
  output_apply_bqsr:
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
      -R $(inputs.indexed_reference_fasta) -I $(inputs.input_bam) -bqsr $(inputs.input_recalibrate_bqsr) -O $(inputs.output_apply_bqsr)

baseCommand: [gatk, ApplyBQSR]
