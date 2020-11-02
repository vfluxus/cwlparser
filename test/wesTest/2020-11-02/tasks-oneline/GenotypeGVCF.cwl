cwlVersion: v1.0

class: CommandLineTool
id: sort_bam
requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement 
    ramMin: 2000

inputs:
  variant_calling: 
    type: File
  indexed_reference_fasta:
    type: File 
    secondaryFiles: [.64.amb, .64.ann, .64.bwt, .64.pac, .64.sa, .64.alt, .dict, .fai]
  output_vcf:
    type: string 
  
outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.vcf'
    secondaryFiles: [^.tbi, ^.vcf.md5]

arguments:
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      -R $(inputs.indexed_reference_fasta) -V $(inputs.variant_calling) -O $(input.output_vcf)

baseCommand: [gatk, GenotypeGVCFs]
