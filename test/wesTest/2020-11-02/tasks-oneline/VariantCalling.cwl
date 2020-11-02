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
  padding_bed:
    type: File
  indexed_reference_fasta:
    type: File 
    secondaryFiles: [.64.amb, .64.ann, .64.bwt, .64.pac, .64.sa, .64.alt, .dict, .fai]
  output_variant_calling:
    type: string 
  
outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.vcf$'
    secondaryFiles: [^.tbi, ^.vcf.md5]

arguments:
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      -R $(inputs.indexed_reference_fasta) -L $(inputs.padding_bed) -I $(inputs.input_bam) -O $(input.output_variant_calling) -ERC GVCF

baseCommand: [gatk, HaplotypeCaller]
